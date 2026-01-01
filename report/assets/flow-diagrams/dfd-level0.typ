// 0层数据流图 - 使用 fletcher 包
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#import fletcher.shapes: diamond, ellipse, rect

#diagram(
  spacing: (18pt, 15pt),
  node-stroke: 1pt,
  edge-stroke: 0.6pt,

  // ===== 外部实体（四角） =====
  node(
    (0, 0),
    text(size: 8pt, weight: "bold")[用户],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#fff9c4"),
    width: 1.5cm,
    height: 0.6cm,
    name: <user>,
  ),
  node(
    (8, 0),
    text(size: 8pt, weight: "bold")[员工],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#fff9c4"),
    width: 1.5cm,
    height: 0.6cm,
    name: <staff>,
  ),
  node(
    (0, 9),
    text(size: 8pt, weight: "bold")[快递员],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#fff9c4"),
    width: 1.5cm,
    height: 0.6cm,
    name: <courier>,
  ),
  node(
    (8, 9),
    text(size: 8pt, weight: "bold")[管理员],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#fff9c4"),
    width: 1.5cm,
    height: 0.6cm,
    name: <admin>,
  ),

  // ===== 处理过程（围绕 D2 分布） =====
  node(
    (2.2, 2.3),
    [#text(size: 7pt, weight: "bold")[P1]\ #text(size: 7pt)[用户管理]],
    stroke: 1pt,
    shape: "circle",
    fill: rgb("#e3f2fd"),
    width: 1.6cm,
    height: 1.6cm,
    name: <p1>,
  ),
  node(
    (5.8, 2.3),
    [#text(size: 7pt, weight: "bold")[P2]\ #text(size: 7pt)[包裹入库]],
    stroke: 1pt,
    shape: "circle",
    fill: rgb("#e3f2fd"),
    width: 1.6cm,
    height: 1.6cm,
    name: <p2>,
  ),
  node(
    (2.2, 5.8),
    [#text(size: 7pt, weight: "bold")[P3]\ #text(size: 7pt)[包裹取件]],
    stroke: 1pt,
    shape: "circle",
    fill: rgb("#e3f2fd"),
    width: 1.6cm,
    height: 1.6cm,
    name: <p3>,
  ),
  node(
    (5.8, 5.8),
    [#text(size: 7pt, weight: "bold")[P4]\ #text(size: 7pt)[滞留处理]],
    stroke: 1pt,
    shape: "circle",
    fill: rgb("#e3f2fd"),
    width: 1.6cm,
    height: 1.6cm,
    name: <p4>,
  ),
  node(
    (4, 8.5),
    [#text(size: 7pt, weight: "bold")[P5]\ #text(size: 7pt)[寄件管理]],
    stroke: 1pt,
    shape: "circle",
    fill: rgb("#e3f2fd"),
    width: 1.6cm,
    height: 1.6cm,
    name: <p5>,
  ),

  // ===== 数据存储（D2 居中） =====
  node(
    (4, 0.8),
    [#text(size: 6pt)[D1] │ #text(size: 6pt)[用户表]],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#f5f5f5"),
    inset: 4pt,
    name: <d1>,
  ),
  node(
    (4, 4.2),
    [#text(size: 6pt)[D2] │ #text(size: 6pt)[包裹表]],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#f5f5f5"),
    inset: 4pt,
    name: <d2>,
  ),
  node(
    (7.2, 4.2),
    [#text(size: 6pt)[D3] │ #text(size: 6pt)[货架表]],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#f5f5f5"),
    inset: 4pt,
    name: <d3>,
  ),
  node(
    (0.8, 4.2),
    [#text(size: 6pt)[D4] │ #text(size: 6pt)[滞留记录]],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#f5f5f5"),
    inset: 4pt,
    name: <d4>,
  ),
  node(
    (4, 7.2),
    [#text(size: 6pt)[D5] │ #text(size: 6pt)[寄件表]],
    stroke: 1pt,
    shape: rect,
    fill: rgb("#f5f5f5"),
    inset: 4pt,
    name: <d5>,
  ),

  // ================= 数据流（按 mermaid 参考语义） =================

  // --- P1 用户管理 ---
  edge(<user>, <p1>, "-|>", label: text(size: 6pt)[注册/登录信息], label-side: left),
  edge(<p1>, <user>, "-|>", label: text(size: 6pt)[认证结果], label-side: right, bend: -15deg),
  edge(<p1>, <d1>, "-|>", label: text(size: 6pt)[读写用户信息], label-side: left, bend: -10deg),
  edge(<d1>, <p1>, "-|>", label: text(size: 6pt)[用户信息], label-side: right, bend: 10deg),

  // --- P2 包裹入库 ---
  edge(<staff>, <p2>, "-|>", label: text(size: 6pt)[入库请求/扫描], label-side: right),
  edge(<p2>, <staff>, "-|>", label: text(size: 6pt)[入库反馈], label-side: left, bend: 15deg),
  edge(<p2>, <d2>, "-|>", label: text(size: 6pt)[录入包裹信息], label-side: right),
  edge(<p2>, <d3>, "-|>", label: text(size: 6pt)[分配/查询货架], label-side: right, bend: 10deg),
  edge(<d3>, <p2>, "-|>", label: text(size: 6pt)[货架信息], label-side: left, bend: -10deg),

  // --- P3 包裹取件 ---
  edge(<user>, <p3>, "-|>", label: text(size: 6pt)[取件码], label-side: left, bend: 25deg),
  edge(<p3>, <user>, "-|>", label: text(size: 6pt)[取件结果], label-side: right, bend: -25deg),
  edge(<p3>, <d2>, "-|>", label: text(size: 6pt)[验证/更新状态], label-side: left),
  edge(<d2>, <p3>, "-|>", label: text(size: 6pt)[包裹状态], label-side: right, bend: 20deg),
  edge(
    <p3>,
    <d3>,
    "-|>",
    stroke: (thickness: 0.6pt, dash: "dashed"),
    label: text(size: 6pt)[释放货架位],
    label-side: right,
    bend: -20deg,
  ),

  // --- P4 滞留处理 ---
  edge(<admin>, <p4>, "-|>", label: text(size: 6pt)[触发检查/系统配置], label-side: right),
  edge(<p4>, <admin>, "-|>", label: text(size: 6pt)[统计报表], label-side: left, bend: -10deg),
  edge(<p4>, <d2>, "-|>", label: text(size: 6pt)[扫描超期包裹], label-side: right),
  edge(<p4>, <d4>, "-|>", label: text(size: 6pt)[写入滞留记录], label-side: left, bend: 10deg),
  edge(<p4>, <user>, "-|>", label: text(size: 6pt)[发送提醒通知], label-side: left, bend: -15deg),

  // --- P5 寄件管理 ---
  edge(<user>, <p5>, "-|>", label: text(size: 6pt)[寄件请求], label-side: left),
  edge(<p5>, <user>, "-|>", label: text(size: 6pt)[运费/单号], label-side: right, bend: 20deg),
  edge(<p5>, <d5>, "-|>", label: text(size: 6pt)[写入寄件记录], label-side: right),
  edge(<p5>, <courier>, "-|>", label: text(size: 6pt)[通知取件], label-side: left, bend: 15deg),
  edge(<courier>, <p5>, "-|>", label: text(size: 6pt)[揽收确认], label-side: right, bend: -15deg),
)
