// 包裹入库流程图 - 使用 fletcher 包
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
    text(size: 9pt)[快递员送达包裹],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <deliver>,
  ),
  node(
    (0, 2),
    text(size: 9pt)[扫描/录入快递单号],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <scan>,
  ),

  // 判断节点
  node(
    (0, 3.2),
    text(size: 8pt)[OCR识别?],
    stroke: 1pt,
    shape: diamond,
    fill: rgb("#fff3e0"),
    inset: 8pt,
    name: <ocr-check>,
  ),

  // 分支节点
  node(
    (2, 3.2),
    text(size: 9pt)[手工补充信息],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <manual>,
  ),
  node(
    (0, 4.9),
    text(size: 9pt)[自动识别收件人信息],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <auto-recognize>,
  ),

  // 后续流程
  node(
    (0, 6),
    text(size: 9pt)[货架分配算法选位],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <assign>,
  ),
  node(
    (0, 7),
    text(size: 9pt)[生成6位取件码],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <generate-code>,
  ),
  node(
    (0, 8),
    text(size: 9pt)[发送取件通知],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <notify>,
  ),

  // 结束
  node(
    (0, 9),
    text(size: 9pt)[结束],
    stroke: 1pt,
    shape: rect,
    corner-radius: 12pt,
    fill: rgb("#e8f5e9"),
    inset: 8pt,
    name: <end>,
  ),

  // 连接线
  edge(<start>, <deliver>, "-|>"),
  edge(<deliver>, <scan>, "-|>"),
  edge(<scan>, <ocr-check>, "-|>"),

  // 判断分支
  edge(<ocr-check>, <manual>, "-|>", label: text(size: 8pt)[否], label-side: left),
  edge(<ocr-check>, <auto-recognize>, "-|>", label: text(size: 8pt)[是], label-side: left),

  // 手工补充汇入主流程
  edge(<manual>, (2, 4.9), stroke: 0.8pt),
  edge((2, 4.9), <auto-recognize>, "-|>"),

  // 后续流程连接
  edge(<auto-recognize>, <assign>, "-|>"),
  edge(<assign>, <generate-code>, "-|>"),
  edge(<generate-code>, <notify>, "-|>"),
  edge(<notify>, <end>, "-|>"),
)
