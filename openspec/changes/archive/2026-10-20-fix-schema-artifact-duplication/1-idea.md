# Idea: show every schema artifact once

As a maintainer using custom OpenSpec schemas (ramona-plan with numbered
outputs like `1-idea.md`) I need dossier to render each schema artifact and
its file exactly once so that the index stays a trustworthy map of the
change.

Today `EnrichChange` merges CLI status artifacts (IDs `idea`, `plan`, …)
with filesystem-discovered artifacts keyed by filename stem (`1-idea`,
`2-plan`, …). Because the two ID spaces never match, every output file is
appended twice — once under the schema artifact and once under a phantom
discovered artifact. Duplicated rows also leak into the artifact tab bar
and make `h`/`l` navigation step over ghost tabs.
