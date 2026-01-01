// 滞留处理流程图 - 使用 fletcher 包
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#import fletcher.shapes: diamond, ellipse, rect

#diagram(
  spacing: (20pt, 15pt),
  node-stroke: 1pt,
  edge-stroke: 0.8pt,

  // 开始
  node(
    (0, 0),
    text(size: 9pt)[定时任务触发],
    stroke: 1pt,
    shape: rect,
    corner-radius: 12pt,
    fill: rgb("#e8f5e9"),
    inset: 8pt,
    name: <start>,
  ),

  // 扫描流程
  node(
    (0, 1),
    text(size: 9pt)[扫描未取件包裹],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <scan>,
  ),

  // 判断节点1
  node(
    (0, 2.2),
    text(size: 8pt)[超过预期时间?],
    stroke: 1pt,
    shape: diamond,
    fill: rgb("#fff3e0"),
    inset: 8pt,
    name: <time-check>,
  ),

  // 未超期分支
  node(
    (3, 2.2),
    text(size: 9pt)[等待下次扫描],
    stroke: 1pt,
    shape: rect,
    corner-radius: 12pt,
    fill: rgb("#e8f5e9"),
    inset: 8pt,
    name: <wait-next>,
  ),

  // 超期处理流程
  node(
    (0, 3.8),
    text(size: 9pt)[标记为滞留状态],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <mark>,
  ),
  node(
    (0, 5),
    text(size: 9pt)[创建滞留记录],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <create-record>,
  ),
  node(
    (0, 6),
    text(size: 9pt)[发送取件提醒],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <remind>,
  ),

  // 判断节点2
  node(
    (0, 7.2),
    text(size: 8pt)[滞留超7天?],
    stroke: 1pt,
    shape: diamond,
    fill: rgb("#fff3e0"),
    inset: 8pt,
    name: <day-check>,
  ),

  // 未超7天分支
  node(
    (3, 7.2),
    text(size: 8pt, "提醒<3次?"),
    stroke: 1pt,
    shape: diamond,
    fill: rgb("#fff3e0"),
    inset: 8pt,
    name: <remind-check>,
  ),

  // 提醒次数判断分支
  node(
    (3, 8.6),
    text(size: 9pt)[等待24小时后],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <wait-24h>,
  ),
  node(
    (5, 7.2),
    text(size: 9pt)[继续等待],
    stroke: 1pt,
    shape: rect,
    corner-radius: 12pt,
    fill: rgb("#e8f5e9"),
    inset: 8pt,
    name: <continue-wait>,
  ),

  // 超7天处理流程
  node(
    (0, 9.8),
    text(size: 9pt)[触发退回流程],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <return>,
  ),
  node(
    (0, 10.9),
    text(size: 9pt)[计算滞留费用],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#e3f2fd"),
    inset: 8pt,
    name: <calc-fee>,
  ),

  // 结束
  node(
    (0, 11.9),
    text(size: 9pt)[流程结束],
    stroke: 1pt,
    shape: rect,
    corner-radius: 12pt,
    fill: rgb("#e8f5e9"),
    inset: 8pt,
    name: <end>,
  ),

  // 主流程连接
  edge(<start>, <scan>, "-|>"),
  edge(<scan>, <time-check>, "-|>"),

  // 时间判断分支
  edge(<time-check>, <wait-next>, "-|>", label: text(size: 8pt)[否], label-side: left),
  edge(<time-check>, <mark>, "-|>", label: text(size: 8pt)[是], label-side: left),

  // 滞留处理流程
  edge(<mark>, <create-record>, "-|>"),
  edge(<create-record>, <remind>, "-|>"),
  edge(<remind>, <day-check>, "-|>"),

  // 天数判断分支
  edge(<day-check>, <remind-check>, "-|>", label: text(size: 8pt)[否], label-side: left),
  edge(<day-check>, <return>, "-|>", label: text(size: 8pt)[是],label-pos: 60%, label-side: left),

  // 提醒次数判断分支
  edge(<remind-check>, <wait-24h>, "-|>", label: text(size: 8pt)[是], label-side: left),
  edge(<remind-check>, <continue-wait>, "-|>", label: text(size: 8pt)[否], label-side: left),

  // 等待24小时后返回发送提醒
  edge(<wait-24h>, (-1.5, 8.6), stroke: 0.8pt),
  edge((-1.5, 8.6), (-1.5, 6), stroke: 0.8pt),
  edge((-1.5, 6), <remind>, "-|>"),

  // 完成流程
  edge(<return>, <calc-fee>, "-|>"),
  edge(<calc-fee>, <end>, "-|>"),
)
