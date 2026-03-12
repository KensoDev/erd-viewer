const COL_H = 24, HDR_H = 32, PAD = 12, MIN_W = 200;

async function main() {
  const schema = await fetch('/schema').then(r => r.json());
  document.getElementById('title').textContent = schema.title;
  document.title = schema.title;

  // ── Layout: simple grid ───────────────────────────────────────────────────
  const COLS = Math.ceil(Math.sqrt(schema.tables.length));
  const GAP_X = 60, GAP_Y = 60;

  const tableMap = {}; // name → { x, y, w, h, cols }
  schema.tables.forEach((t, i) => {
    const w = Math.max(MIN_W, Math.max(...t.columns.map(c => c.name.length * 7.5 + c.type.length * 6.5 + 40)) + PAD * 2);
    const h = HDR_H + t.columns.length * COL_H;
    tableMap[t.name] = { x: 0, y: 0, w, h, cols: t.columns, name: t.name };
  });

  // Grid placement
  let colWidths = [], rowHeights = [];
  schema.tables.forEach((t, i) => {
    const col = i % COLS, row = Math.floor(i / COLS);
    const { w, h } = tableMap[t.name];
    colWidths[col]  = Math.max(colWidths[col]  || 0, w);
    rowHeights[row] = Math.max(rowHeights[row] || 0, h);
  });
  schema.tables.forEach((t, i) => {
    const col = i % COLS, row = Math.floor(i / COLS);
    let x = GAP_X, y = GAP_Y + 48;
    for (let c = 0; c < col; c++) x += (colWidths[c] || 0) + GAP_X;
    for (let r = 0; r < row; r++) y += (rowHeights[r] || 0) + GAP_Y;
    tableMap[t.name].x = x;
    tableMap[t.name].y = y;
  });

  // ── D3 setup ──────────────────────────────────────────────────────────────
  const svg    = d3.select('#canvas');
  const root   = d3.select('#root');
  const eGroup = d3.select('#edges');
  const nGroup = d3.select('#nodes');

  const zoom = d3.zoom().scaleExtent([0.1, 3]).on('zoom', e => root.attr('transform', e.transform));
  svg.call(zoom);

  // ── Draw edges ────────────────────────────────────────────────────────────
  function edgePath(fk) {
    const s = tableMap[fk.fromTable], t = tableMap[fk.toTable];
    if (!s || !t) return '';
    const sx = s.x + s.w, sy = s.y + HDR_H + s.cols.findIndex(c => c.name === fk.fromCol) * COL_H + COL_H / 2;
    const tx = t.x,       ty = t.y + HDR_H + t.cols.findIndex(c => c.name === fk.toCol)   * COL_H + COL_H / 2;
    const mx = (sx + tx) / 2;
    return `M${sx},${sy} C${mx},${sy} ${mx},${ty} ${tx},${ty}`;
  }

  const edges = eGroup.selectAll('g.fk').data(schema.fks).enter().append('g').attr('class', 'fk');
  edges.append('path').attr('class', 'fk-line').attr('d', edgePath).attr('marker-end', 'url(#arrow-fk)');
  edges.append('text').attr('class', 'fk-label').each(function(fk) {
    const s = tableMap[fk.fromTable], t = tableMap[fk.toTable];
    if (!s || !t) return;
    const sx = s.x + s.w, sy = s.y + HDR_H + s.cols.findIndex(c => c.name === fk.fromCol) * COL_H + COL_H / 2;
    const tx = t.x,       ty = t.y + HDR_H + t.cols.findIndex(c => c.name === fk.toCol)   * COL_H + COL_H / 2;
    d3.select(this).attr('x', (sx + tx) / 2).attr('y', (sy + ty) / 2 - 4)
      .attr('text-anchor', 'middle').text(`${fk.fromCol} → ${fk.toCol}`);
  });

  function refreshEdges() {
    edges.select('path').attr('d', edgePath);
    edges.select('text').each(function(fk) {
      const s = tableMap[fk.fromTable], t = tableMap[fk.toTable];
      if (!s || !t) return;
      const sx = s.x + s.w, sy = s.y + HDR_H + s.cols.findIndex(c => c.name === fk.fromCol) * COL_H + COL_H / 2;
      const tx = t.x,       ty = t.y + HDR_H + t.cols.findIndex(c => c.name === fk.toCol)   * COL_H + COL_H / 2;
      d3.select(this).attr('x', (sx + tx) / 2).attr('y', (sy + ty) / 2 - 4);
    });
  }

  // ── Draw nodes ────────────────────────────────────────────────────────────
  const tooltip = document.getElementById('tooltip');

  const nodes = nGroup.selectAll('g.table-node')
    .data(schema.tables).enter()
    .append('g').attr('class', 'table-node')
    .attr('transform', t => `translate(${tableMap[t.name].x},${tableMap[t.name].y})`)
    .call(d3.drag()
      .on('start', function() { d3.select(this).raise(); })
      .on('drag', function(event, t) {
        const m = tableMap[t.name];
        m.x += event.dx; m.y += event.dy;
        d3.select(this).attr('transform', `translate(${m.x},${m.y})`);
        refreshEdges();
      })
    )
    .on('mouseover', function(event, t) {
      const m = tableMap[t.name];
      // Highlight connected edges + nodes
      const related = new Set([t.name]);
      edges.each(fk => { if (fk.fromTable === t.name || fk.toTable === t.name) { related.add(fk.fromTable); related.add(fk.toTable); }});
      nodes.classed('dimmed', n => !related.has(n.name)).classed('highlighted', n => n.name === t.name);
      edges.select('path')
        .classed('highlighted', fk => fk.fromTable === t.name || fk.toTable === t.name)
        .attr('marker-end', fk => (fk.fromTable === t.name || fk.toTable === t.name) ? 'url(#arrow-fk-hl)' : 'url(#arrow-fk)');
      // Tooltip
      const fksFrom = schema.fks.filter(fk => fk.fromTable === t.name);
      const fksTo   = schema.fks.filter(fk => fk.toTable   === t.name);
      let html = `<div class="tt-title">${t.name}</div>`;
      html += `<div class="tt-row">${t.columns.length} columns</div>`;
      if (fksFrom.length) html += `<div class="tt-row">→ references: ${fksFrom.map(f => f.toTable).join(', ')}</div>`;
      if (fksTo.length)   html += `<div class="tt-row">← referenced by: ${fksTo.map(f => f.fromTable).join(', ')}</div>`;
      tooltip.innerHTML = html;
      tooltip.style.opacity = 1;
    })
    .on('mousemove', event => {
      tooltip.style.left = (event.clientX + 14) + 'px';
      tooltip.style.top  = (event.clientY + 14) + 'px';
    })
    .on('mouseout', function() {
      nodes.classed('dimmed', false).classed('highlighted', false);
      edges.select('path').classed('highlighted', false).attr('marker-end', 'url(#arrow-fk)');
      tooltip.style.opacity = 0;
    });

  nodes.each(function(t) {
    const m = tableMap[t.name];
    const g = d3.select(this);
    // Body background
    g.append('rect').attr('class', 'tbl-body').attr('width', m.w).attr('height', m.h).attr('rx', 6);
    // Header
    g.append('rect').attr('class', 'tbl-header').attr('width', m.w).attr('height', HDR_H).attr('rx', 6);
    g.append('rect').attr('fill', '#3b82f6').attr('width', m.w).attr('height', HDR_H / 2).attr('y', HDR_H / 2); // square bottom corners
    g.append('text').attr('class', 'tbl-header-text').attr('x', m.w / 2).attr('y', 21).attr('text-anchor', 'middle').text(t.name);
    // Columns
    t.columns.forEach((col, i) => {
      const y = HDR_H + i * COL_H;
      g.append('rect').attr('class', col.IsPK ? 'col-row-pk' : (i % 2 === 0 ? 'col-row' : ''))
        .attr('x', 1).attr('y', y).attr('width', m.w - 2).attr('height', COL_H)
        .attr('fill', col.IsPK ? '#fef3c7' : (i % 2 === 0 ? '#1e293b' : '#243044'));
      const icon = col.IsPK ? '🔑 ' : '';
      g.append('text').attr('class', col.IsPK ? 'col-name-pk' : 'col-name')
        .attr('x', PAD).attr('y', y + 16).text(icon + col.name);
      g.append('text').attr('class', 'col-type')
        .attr('x', m.w - PAD).attr('y', y + 16).attr('text-anchor', 'end')
        .text(col.type + (col.Nullable ? '' : ' *'));
    });
  });

  // ── Toolbar ───────────────────────────────────────────────────────────────
  function fitView() {
    const w = window.innerWidth, h = window.innerHeight;
    let minX = Infinity, minY = Infinity, maxX = -Infinity, maxY = -Infinity;
    Object.values(tableMap).forEach(m => {
      minX = Math.min(minX, m.x); minY = Math.min(minY, m.y);
      maxX = Math.max(maxX, m.x + m.w); maxY = Math.max(maxY, m.y + m.h);
    });
    const scale = Math.min(0.9, Math.min((w - 40) / (maxX - minX), (h - 80) / (maxY - minY)));
    svg.transition().duration(400).call(zoom.transform, d3.zoomIdentity.translate(w / 2, h / 2).scale(scale).translate(-(minX + maxX) / 2, -(minY + maxY) / 2));
  }

  document.getElementById('btn-fit').addEventListener('click', fitView);
  document.getElementById('btn-reset').addEventListener('click', () => {
    schema.tables.forEach((t, i) => {
      const col = i % COLS, row = Math.floor(i / COLS);
      let x = GAP_X, y = GAP_Y + 48;
      for (let c = 0; c < col; c++) x += (colWidths[c] || 0) + GAP_X;
      for (let r = 0; r < row; r++) y += (rowHeights[r] || 0) + GAP_Y;
      tableMap[t.name].x = x; tableMap[t.name].y = y;
    });
    nodes.attr('transform', t => `translate(${tableMap[t.name].x},${tableMap[t.name].y})`);
    refreshEdges();
    fitView();
  });

  document.getElementById('search').addEventListener('input', e => {
    const q = e.target.value.trim().toLowerCase();
    if (!q) { nodes.classed('dimmed', false); return; }
    nodes.classed('dimmed', t => !t.name.toLowerCase().includes(q));
  });

  fitView();
}

main();
