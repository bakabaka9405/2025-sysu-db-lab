// 包裹取件流程图 - 使用 fletcher 包
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#import fletcher.shapes: diamond, ellipse, rect

#diagram(
  spacing: (20pt, 15pt),
  node-stroke: 1pt,
  edge-stroke: 0.8pt,

  // 开始
  node(
    (0, 0),
    text(size: 9pt)[开始],
    stroke: 1pt,
    shape: rect,
    corner-radius: 12pt,
    fill: rgb("#e8f5e9"),
    inset: 8pt,
    name: <start>,
  ),

  // 流程节点
  node(
    (0, 1),
    text(size: 9pt)[用户提供取件码],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <provide-code>,
  ),

  // 判断节点1
  node(
    (0, 2.4),
    text(size: 8pt)[验证通过?],
    stroke: 1pt,
    shape: diamond,
    fill: rgb("#fff3e0"),
    inset: 8pt,
    name: <verify-check>,
  ),

  // 验证失败分支
  node(
    (1.8, 2.4),
    text(size: 9pt)[提示验证失败],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <verify-fail>,
  ),

  // 验证成功后续流程
  node(
    (0, 4),
    text(size: 9pt)[显示包裹与货架信息],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <show-info>,
  ),
  node(
    (0, 5),
    text(size: 9pt)[工作人员取出包裹],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <fetch>,
  ),

  // 判断节点2
  node(
    (0, 6.4),
    text(size: 8pt)[核对无误?],
    stroke: 1pt,
    shape: diamond,
    fill: rgb("#fff3e0"),
    inset: 8pt,
    name: <confirm-check>,
  ),

  // 核对失败分支
  node(
    (1, 6.4),
    text(size: 9pt)[重新核对信息],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <recheck>,
  ),

  // 核对成功后续流程
  node(
    (0, 8),
    text(size: 9pt)[用户签收确认],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <sign>,
  ),
  node(
    (0, 9),
    text(size: 9pt)[更新状态/释放货架],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <update>,
  ),

  // 结束
  node(
    (0, 10),
    text(size: 9pt)[结束],
    stroke: 1pt,
    shape: rect,
    corner-radius: 12pt,
    fill: rgb("#e8f5e9"),
    inset: 8pt,
    name: <end>,
  ),

  // 主流程连接
  edge(<start>, <provide-code>, "-|>"),
  edge(<provide-code>, <verify-check>, "-|>"),

  // 验证判断分支
  edge(<verify-check>, <verify-fail>, "-|>", label: text(size: 8pt)[否], label-side: left),
  edge(<verify-check>, <show-info>, "-|>", label: text(size: 8pt)[是], label-side: left),

  // 验证失败返回结束
  edge(<verify-fail>, (1.8, 10), stroke: 0.8pt),
  edge((1.8, 10), <end>, "-|>"),

  // 继续主流程
  edge(<show-info>, <fetch>, "-|>"),
  edge(<fetch>, <confirm-check>, "-|>"),

  // 核对判断分支
  edge(<confirm-check>, <recheck>, "-|>", label: text(size: 8pt)[否], label-side: left),
  edge(<confirm-check>, <sign>, "-|>", label: text(size: 8pt)[是], label-side: left),

  // 重新核对返回取出包裹
  edge(<recheck>, (1, 5), stroke: 0.8pt),
  edge((1, 5), <fetch>, "-|>"),

  // 完成流程
  edge(<sign>, <update>, "-|>"),
  edge(<update>, <end>, "-|>"),
)
