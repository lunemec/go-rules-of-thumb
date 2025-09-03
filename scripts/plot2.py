import pandas as pd
import matplotlib.pyplot as plt
import re
import argparse
from io import StringIO
import seaborn as sns

# --- CLI ARGUMENTS ---
parser = argparse.ArgumentParser(
    description="Plot ns/op from benchstat CSV with multi-param support"
)
parser.add_argument("csvfile", help="Path to benchstat CSV file")
parser.add_argument("output", help="Output PNG file")
args = parser.parse_args()

# --- LOAD CSV LINES ---
with open(args.csvfile) as f:
    lines = f.readlines()

rows = []
file_labels = []

# --- FIND sec/op SECTION ---
for i, line in enumerate(lines):
    if "sec/op" in line and "CI" in line:
        file_line = lines[i - 1]
        file_parts = [p.strip() for p in file_line.split(",") if p.strip()]
        file_labels = [label.removesuffix(".txt") for label in file_parts]

        # Collect sec/op rows
        j = i + 1
        while j < len(lines):
            if (
                lines[j].strip() == ""
                or lines[j].startswith(",")
                or lines[j].lower().startswith("geomean")
            ):
                break
            rows.append(lines[j])
            j += 1
        break

# --- PARSE DATA ROWS ---
csv_data = "".join(rows)
df_raw = pd.read_csv(StringIO(csv_data), header=None)


# --- EXTRACT PARAMETERS FROM BENCHMARK NAME ---
def extract_params(name):
    # Split once at '/', take the right side
    right = name.split("/", 1)[1]
    # Split once at '-', take the left side (before thread count suffix)
    param_part = right.split("-", 1)[0]
    # Now split into key=value pairs
    params = {}
    for part in param_part.split("_"):
        if "=" in part:
            key, value = part.split("=", 1)
            try:
                params[key] = int(value)
            except ValueError:
                params[key] = value  # fallback for non-integer values
    return params


# Extract parameters for each row as a dictionary
param_dicts = df_raw[0].apply(extract_params)
# Convert list of dicts to DataFrame and concatenate
param_df = pd.DataFrame(param_dicts.tolist())
# Merge parameter columns into df_raw
df_raw = pd.concat([df_raw, param_df], axis=1)

# --- FIND sec/op COLUMN INDICES ---
secop_header = [h.strip() for h in lines[i].split(",")]
secop_indices = [idx for idx, val in enumerate(secop_header) if val == "sec/op"]

if len(secop_indices) != len(file_labels):
    raise ValueError(
        f"Mismatch: Found {len(secop_indices)} 'sec/op' columns but {len(file_labels)} file labels"
    )

# --- CONVERT TO ns/op (Generic) ---
records = []
for idx, label in zip(secop_indices, file_labels):
    for _, row in df_raw.iterrows():
        try:
            ns_per_op = float(row[idx]) * 1e9
        except (ValueError, TypeError):
            continue  # skip if conversion fails

        # Extract all params (everything beyond fixed CSV columns)
        row_params = {
            k: row[k]
            for k in df_raw.columns
            if k not in df_raw.columns[: max(secop_indices) + 1]
        }

        records.append(
            {
                "implementation": label,
                "ns_per_op": ns_per_op,
                **row_params,  # dynamically unpack all parameters
            }
        )

df = pd.DataFrame(records)

# Fixed X axis
x_col = "size"
y_col = "ns_per_op"

# Choose hue_col: prefer 'needles' if it exists, otherwise use first other param
# if "iterations" in df.columns:
#    hue_col = "iterations"
#    df[hue_col] = pd.to_numeric(df[hue_col], errors="coerce")
if "subset" in df.columns:
    hue_col = "subset"
    df[hue_col] = pd.to_numeric(df[hue_col], errors="coerce")
else:
    hue_col = "implementation"

# Ensure numeric conversion
df[x_col] = pd.to_numeric(df[x_col], errors="coerce")
df[y_col] = pd.to_numeric(df[y_col], errors="coerce")

plt.figure(figsize=(12, 6))
plot_args = dict(
    data=df,
    x=x_col,
    y=y_col,
    hue=hue_col,
    style="implementation",
    col="subset",
    col_wrap=1,
    markers=True,
    dashes=False,
    palette="deep",
    kind="line",
)
print(plot_args)
print(df)
sns.set_theme(style="whitegrid")
sns.relplot(**plot_args)
plt.xscale("log")
plt.yscale("log")
plt.xlabel("Size")
plt.ylabel("ns/op")
plt.tight_layout()
plt.savefig(args.output)
print(f"Saved graph to {args.output}")
