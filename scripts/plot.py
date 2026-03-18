import argparse
import csv
import json
import os
import re
import tempfile
from io import StringIO
from pathlib import Path

cache_root = Path(tempfile.gettempdir()) / "go-rules-of-thumb-cache"
cache_root.mkdir(parents=True, exist_ok=True)
matplotlib_cache = cache_root / "matplotlib"
xdg_cache = cache_root / "xdg"
matplotlib_cache.mkdir(parents=True, exist_ok=True)
xdg_cache.mkdir(parents=True, exist_ok=True)
(xdg_cache / "fontconfig").mkdir(parents=True, exist_ok=True)
os.environ.setdefault("MPLCONFIGDIR", str(matplotlib_cache))
os.environ.setdefault("XDG_CACHE_HOME", str(xdg_cache))

import matplotlib

matplotlib.use("Agg")
import matplotlib.pyplot as plt
import pandas as pd
import seaborn as sns
from matplotlib.colors import ListedColormap, LogNorm
from matplotlib.patches import Patch, Rectangle


PLOT_SPECS_PATH = Path(__file__).with_name("plots.json")


def parse_args():
    parser = argparse.ArgumentParser(
        description="Plot ns/op from benchstat CSV using benchmark-specific specs",
    )
    parser.add_argument("csvfile", help="Path to a benchstat CSV file")
    parser.add_argument("output", help="Primary output PNG path")
    return parser.parse_args()


def read_plot_specs():
    with PLOT_SPECS_PATH.open(encoding="utf-8") as handle:
        return json.load(handle)


def read_csv_cells(line):
    return next(csv.reader([line]))


def benchmark_name_from_csv(csv_path):
    return Path(csv_path).stem


def benchmark_title(benchmark_name):
    title = benchmark_name.removeprefix("Benchmark")
    return re.sub(r"(?<!^)(?=[A-Z])", " ", title)


def load_spec(benchmark_name):
    specs = read_plot_specs()
    try:
        return specs[benchmark_name]
    except KeyError as exc:
        known = ", ".join(sorted(specs))
        raise SystemExit(
            f"Missing plot config for {benchmark_name}. Known benchmarks: {known}"
        ) from exc


def extract_params(raw_name):
    if "/" in raw_name:
        param_part = raw_name.split("/", 1)[1]
    else:
        param_part = raw_name

    param_part = re.sub(r"-\d+$", "", param_part)
    params = {}
    for chunk in param_part.split("_"):
        if "=" not in chunk:
            continue
        key, value = chunk.split("=", 1)
        try:
            params[key] = int(value)
        except ValueError:
            params[key] = value
    return params


def normalize_implementation_name(label, benchmark_name):
    prefix = f"{benchmark_name}-"
    if label.startswith(prefix):
        return label[len(prefix) :]
    return label


def parse_benchstat(csv_path, benchmark_name):
    lines = Path(csv_path).read_text(encoding="utf-8").splitlines()

    rows = []
    file_labels = []
    secop_indices = []

    for idx, line in enumerate(lines):
        header_cells = [cell.strip() for cell in read_csv_cells(line)]
        if "sec/op" not in header_cells or "CI" not in header_cells:
            continue

        if idx == 0:
            raise SystemExit(f"Unexpected benchstat CSV format in {csv_path}")

        file_cells = [cell.strip() for cell in read_csv_cells(lines[idx - 1])]
        file_labels = [cell.removesuffix(".txt") for cell in file_cells if cell]
        secop_indices = [
            column_index
            for column_index, value in enumerate(header_cells)
            if value == "sec/op"
        ]

        row_index = idx + 1
        while row_index < len(lines):
            row = lines[row_index].strip()
            if not row or row.startswith(",") or row.lower().startswith("geomean"):
                break
            rows.append(lines[row_index])
            row_index += 1
        break

    if not rows:
        raise SystemExit(f"No sec/op section found in {csv_path}")

    if len(secop_indices) != len(file_labels):
        raise SystemExit(
            f"Mismatch in {csv_path}: found {len(secop_indices)} sec/op columns "
            f"but {len(file_labels)} benchmark files"
        )

    df_raw = pd.read_csv(StringIO("\n".join(rows)), header=None)
    params = pd.DataFrame(df_raw[0].map(extract_params).tolist())

    records = []
    for secop_index, label in zip(secop_indices, file_labels):
        implementation = normalize_implementation_name(label, benchmark_name)
        for row_idx, row in df_raw.iterrows():
            try:
                ns_per_op = float(row[secop_index]) * 1e9
            except (TypeError, ValueError):
                continue

            record = {
                "implementation": implementation,
                "ns_per_op": ns_per_op,
            }
            for column in params.columns:
                record[column] = params.at[row_idx, column]
            records.append(record)

    df = pd.DataFrame(records)
    if df.empty:
        raise SystemExit(f"No sec/op values could be parsed from {csv_path}")

    for column in params.columns:
        converted = pd.to_numeric(df[column], errors="coerce")
        if not converted.isna().all():
            df[column] = converted

    return df


def validate_spec(df, spec, benchmark_name):
    for required_key in ("x_param", "summary_kind", "detail_kind", "axis_labels"):
        if required_key not in spec:
            raise SystemExit(f"{benchmark_name}: missing required spec key {required_key}")

    varying_params = [
        column
        for column in df.columns
        if column not in ("implementation", "ns_per_op") and df[column].nunique() > 1
    ]

    if not varying_params:
        raise SystemExit(f"{benchmark_name}: no varying benchmark parameters found")
    if len(varying_params) > 2:
        joined = ", ".join(sorted(varying_params))
        raise SystemExit(f"{benchmark_name}: unsupported benchmark dimensions: {joined}")

    expected = [spec["x_param"]]
    if spec.get("y_param") is not None:
        expected.append(spec["y_param"])

    if set(varying_params) != set(expected):
        actual = ", ".join(sorted(varying_params))
        wanted = ", ".join(expected)
        raise SystemExit(
            f"{benchmark_name}: spec expects varying params [{wanted}] but CSV contains [{actual}]"
        )

    if spec["summary_kind"] not in {"line", "winner_heatmap", "value_heatmap"}:
        raise SystemExit(
            f"{benchmark_name}: unsupported summary_kind {spec['summary_kind']}"
        )

    if spec["detail_kind"] not in {None, "facet_lines"}:
        raise SystemExit(
            f"{benchmark_name}: unsupported detail_kind {spec['detail_kind']}"
        )

    return varying_params


def axis_label(spec, axis_name, fallback):
    labels = spec.get("axis_labels", {})
    return labels.get(axis_name, fallback)


def param_label(spec, param_name):
    return spec.get("param_labels", {}).get(param_name, param_name.replace("_", " "))


def metric_label(spec):
    return axis_label(spec, "metric", "ns/op")


def implementation_label(spec, implementation):
    return spec.get("implementation_labels", {}).get(implementation, implementation)


def ordered_implementations(df, spec):
    configured = spec.get("implementation_order", [])
    found = list(df["implementation"].dropna().unique())
    ordered = [implementation for implementation in configured if implementation in found]
    extras = sorted(
        implementation for implementation in found if implementation not in ordered
    )
    return ordered + extras


def with_display_labels(df, spec):
    display_df = df.copy()
    display_df["implementation_label"] = display_df["implementation"].map(
        lambda implementation: implementation_label(spec, implementation)
    )
    return display_df


def save_figure(figure, output_path):
    output_path.parent.mkdir(parents=True, exist_ok=True)
    figure.savefig(output_path, dpi=180, bbox_inches="tight")
    plt.close(figure)


def plot_line_summary(df, spec, benchmark_name, output_path):
    x_param = spec["x_param"]
    implementation_order = ordered_implementations(df, spec)
    display_df = with_display_labels(df, spec)
    display_order = [
        implementation_label(spec, implementation) for implementation in implementation_order
    ]

    fig, ax = plt.subplots(figsize=(12, 6))
    sns.lineplot(
        data=display_df.sort_values([x_param, "implementation"]),
        x=x_param,
        y="ns_per_op",
        hue="implementation_label",
        style="implementation_label",
        hue_order=display_order,
        style_order=display_order,
        markers=True,
        dashes=False,
        ax=ax,
    )

    ax.set_xscale("log")
    ax.set_yscale("log")
    ax.set_xlabel(axis_label(spec, "x", param_label(spec, x_param)))
    ax.set_ylabel(metric_label(spec))
    ax.legend(title="Implementation")
    fig.tight_layout()
    save_figure(fig, output_path)


def format_ns(value):
    if pd.isna(value):
        return ""
    if value < 1_000:
        return f"{round(value):.0f}ns"
    if value < 1_000_000:
        return f"{round(value / 1_000):.0f}µs"
    if value < 1_000_000_000:
        return f"{round(value / 1_000_000):.0f}ms"
    return f"{round(value / 1_000_000_000):.0f}s"


def plot_value_heatmap(df, spec, benchmark_name, output_path):
    x_param = spec["x_param"]
    implementation_order = ordered_implementations(df, spec)
    display_order = [
        implementation_label(spec, implementation) for implementation in implementation_order
    ]
    display_df = with_display_labels(df, spec)
    x_values = sorted(df[x_param].dropna().unique())

    pivot = (
        display_df.pivot(index="implementation_label", columns=x_param, values="ns_per_op")
        .reindex(index=display_order, columns=x_values)
    )
    annotations = pivot.copy().astype(object)
    fastest_cells = set()

    for col_idx, x_value in enumerate(x_values):
        column = pivot[x_value]
        fastest = column.idxmin()
        for row_idx, implementation in enumerate(display_order):
            label = format_ns(pivot.at[implementation, x_value])
            annotations.at[implementation, x_value] = label
            if implementation == fastest and label:
                fastest_cells.add((row_idx, col_idx))

    fig_width = max(10, len(x_values) * 0.95)
    fig_height = max(3.8, len(display_order) * 0.9)
    fig, ax = plt.subplots(figsize=(fig_width, fig_height))

    heatmap = sns.heatmap(
        pivot,
        cmap=sns.blend_palette(["#2a9d8f", "#e9c46a", "#e76f51"], as_cmap=True),
        norm=LogNorm(vmin=float(pivot.min().min()), vmax=float(pivot.max().max())),
        mask=pivot.isna(),
        annot=annotations,
        fmt="",
        linewidths=0.5,
        linecolor="white",
        cbar_kws={"label": metric_label(spec)},
        ax=ax,
    )

    heatmap.set_xlabel(axis_label(spec, "x", param_label(spec, x_param)))
    heatmap.set_ylabel("Implementation\n(one row each)", fontweight="bold")
    heatmap.set_xticklabels([str(value) for value in x_values], rotation=45, ha="right")
    heatmap.set_yticklabels(display_order, rotation=0)
    ax.tick_params(axis="y", labelsize=11)
    for tick in ax.get_yticklabels():
        tick.set_fontweight("bold")

    for text in ax.texts:
        x, y = text.get_position()
        row_idx = round(y - 0.5)
        col_idx = round(x - 0.5)
        if (row_idx, col_idx) in fastest_cells:
            text.set_fontweight("bold")

    for row_idx, col_idx in fastest_cells:
        ax.add_patch(
            Rectangle(
                (col_idx, row_idx),
                1,
                1,
                fill=False,
                edgecolor="#111111",
                linewidth=2.5,
            )
        )

    ax.text(
        0,
        1.08,
        "Each row is one implementation. Bold text and a border mark the fastest value in each column.",
        transform=ax.transAxes,
        ha="left",
        va="bottom",
        fontsize=10,
        fontweight="bold",
    )

    fig.tight_layout()
    save_figure(fig, output_path)


def plot_winner_heatmap(df, spec, benchmark_name, output_path):
    x_param = spec["x_param"]
    y_param = spec["y_param"]
    implementation_order = ordered_implementations(df, spec)
    display_order = [
        implementation_label(spec, implementation) for implementation in implementation_order
    ]

    winners = df.loc[df.groupby([y_param, x_param])["ns_per_op"].idxmin()].copy()
    winners["implementation_label"] = winners["implementation"].map(
        lambda implementation: implementation_label(spec, implementation)
    )
    x_values = sorted(df[x_param].dropna().unique())
    y_values = sorted(df[y_param].dropna().unique(), reverse=True)

    pivot = (
        winners.pivot(index=y_param, columns=x_param, values="implementation_label")
        .reindex(index=y_values, columns=x_values)
    )

    codes = {implementation: idx for idx, implementation in enumerate(display_order)}
    heatmap_values = pivot.apply(lambda column: column.map(codes)).astype(float)
    palette = sns.color_palette("colorblind", n_colors=len(display_order))

    fig_width = max(10, len(x_values) * 0.9)
    fig_height = max(6, len(y_values) * 0.45)
    fig, ax = plt.subplots(figsize=(fig_width, fig_height))
    sns.heatmap(
        heatmap_values,
        cmap=ListedColormap(palette),
        mask=heatmap_values.isna(),
        cbar=False,
        linewidths=0.5,
        linecolor="white",
        vmin=-0.5,
        vmax=len(display_order) - 0.5,
        ax=ax,
    )

    ax.set_xlabel(axis_label(spec, "x", param_label(spec, x_param)))
    ax.set_ylabel(axis_label(spec, "y", param_label(spec, y_param)))
    ax.set_xticklabels([str(value) for value in x_values], rotation=45, ha="right")
    ax.set_yticklabels([str(value) for value in y_values], rotation=0)

    legend_handles = [
        Patch(facecolor=palette[idx], label=implementation)
        for idx, implementation in enumerate(display_order)
    ]
    ax.legend(
        handles=legend_handles,
        title="Fastest implementation",
        bbox_to_anchor=(1.02, 1),
        loc="upper left",
        borderaxespad=0,
    )

    fig.tight_layout()
    save_figure(fig, output_path)


def plot_facet_lines(df, spec, benchmark_name, output_path):
    x_param = spec["x_param"]
    y_param = spec["y_param"]
    implementation_order = ordered_implementations(df, spec)
    display_df = with_display_labels(df, spec)
    display_order = [
        implementation_label(spec, implementation) for implementation in implementation_order
    ]
    facet_values = sorted(df[y_param].dropna().unique())

    chart = sns.relplot(
        data=display_df.sort_values([y_param, x_param, "implementation"]),
        x=x_param,
        y="ns_per_op",
        hue="implementation_label",
        style="implementation_label",
        hue_order=display_order,
        style_order=display_order,
        col=y_param,
        col_order=facet_values,
        col_wrap=min(3, len(facet_values)),
        kind="line",
        markers=True,
        dashes=False,
        facet_kws={"sharex": True, "sharey": True},
        height=3.1,
        aspect=1.2,
    )

    for axis in chart.axes.flat:
        axis.set_xscale("log")
        axis.set_yscale("log")

    chart.set_axis_labels(
        axis_label(spec, "x", param_label(spec, x_param)),
        metric_label(spec),
    )
    chart.set_titles(
        f"{axis_label(spec, 'y', param_label(spec, y_param))} = {{col_name}}"
    )
    if chart._legend is not None:
        chart._legend.set_title("Implementation")
    chart.savefig(output_path, dpi=180, bbox_inches="tight")
    plt.close(chart.fig)


def detail_output_path(output_path):
    return output_path.with_name(f"{output_path.stem}-detail{output_path.suffix}")


def main():
    args = parse_args()
    benchmark_name = benchmark_name_from_csv(args.csvfile)
    spec = load_spec(benchmark_name)
    df = parse_benchstat(args.csvfile, benchmark_name)
    validate_spec(df, spec, benchmark_name)

    sns.set_theme(style="whitegrid")
    output_path = Path(args.output)

    if spec["summary_kind"] == "line":
        plot_line_summary(df, spec, benchmark_name, output_path)
    elif spec["summary_kind"] == "winner_heatmap":
        plot_winner_heatmap(df, spec, benchmark_name, output_path)
    else:
        plot_value_heatmap(df, spec, benchmark_name, output_path)

    print(f"Wrote {output_path}")

    if spec.get("detail_kind") == "facet_lines":
        detail_path = detail_output_path(output_path)
        plot_facet_lines(df, spec, benchmark_name, detail_path)
        print(f"Wrote {detail_path}")


if __name__ == "__main__":
    main()
