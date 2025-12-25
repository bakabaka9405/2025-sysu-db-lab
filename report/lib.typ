#let data-table(
  caption: none,
  columns: (auto, 1fr),
  align: left,
  headers: (),
  ..rows,
) = {
  let content = table(
    columns: columns,
    align: align,
    stroke: 0.5pt,
    inset: 8pt,
    ..if headers.len() > 0 { (table.header(..headers.map(h => [*#h*])),) } else { () },
    ..rows.pos().flatten()
  )
  if caption != none {
    figure(content, caption: caption)
  } else {
    content
  }
}

#let soft-header(content) = {
  par(first-line-indent: 0em, text(content, weight: "bold"))
}
