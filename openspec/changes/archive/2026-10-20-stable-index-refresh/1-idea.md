# Idea: an index that never flaps while work happens underneath it

As a developer keeping dossier open while an agent (or my editor) edits
OpenSpec artifacts I need the index hierarchy to update in place with stable
row identities so that the screen never flashes, rows never rename, and my
cursor never jumps.

Today every disk change discards the CLI-enriched change and swaps in the
filesystem-discovered view whose artifact IDs come from filename stems
(`1-idea`, `2-plan`, …), then asynchronously re-fetches `openspec status`
and swaps back to schema IDs (`idea`, `plan`, …). Each edit therefore
produces two full index rebuilds, two row-identity renames, and a lost
cursor — the "weird refresh". A scratch simulation (`ramona-plan` change,
edit while open) measured: cursor identity flapping
`artifact:idea → artifact:1-idea → artifact:idea` within one 500 ms tick,
with rebuilds at each step; steady state (no edits) is perfectly stable.
