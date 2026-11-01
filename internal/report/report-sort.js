(() => {
  const container = document.getElementById('source-rows');
  const status = document.getElementById('sort-status');
  const buttons = Array.from(document.querySelectorAll('[data-sort]'));
  if (!container || !status) return;
  const rows = Array.from(container.children).map((element, index) => ({ element, index }));
  let column = 'findings';
  let direction = 'desc';
  const compare = (a, b) => a < b ? -1 : a > b ? 1 : 0;
  const pathOrder = (a, b, order) => {
    const x = a.element.dataset;
    const y = b.element.dataset;
    if (Boolean(x.path) !== Boolean(y.path)) return x.path ? -1 : 1;
    const path = compare(x.path, y.path);
    return path * order || compare(x.method, y.method) || a.index - b.index;
  };
  for (const button of buttons) {
    button.disabled = false;
    button.addEventListener('click', () => {
      const next = button.dataset.sort;
      direction = next === column ? (direction === 'asc' ? 'desc' : 'asc') : (next === 'findings' ? 'desc' : 'asc');
      column = next;
      const order = direction === 'asc' ? 1 : -1;
      rows.sort((a, b) => column === 'findings'
        ? ['findings', 'review', 'protected'].reduce((result, key) => result || (Number(a.element.dataset[key]) - Number(b.element.dataset[key])) * order, 0) || pathOrder(a, b, 1)
        : pathOrder(a, b, order));
      for (const row of rows) container.appendChild(row.element);
      for (const control of buttons) {
        const active = control.dataset.sort === column;
        control.setAttribute('aria-pressed', String(active));
        control.querySelector('.sort-direction').textContent = active ? (direction === 'asc' ? ' ↑' : ' ↓') : ' ↕';
      }
      status.textContent = column === 'findings'
        ? `Sorted by Findings, then Review, then Protected (${direction === 'desc' ? 'highest' : 'lowest'} first), then route path A–Z and method.`
        : `Sorted by route path ${direction === 'asc' ? 'A–Z' : 'Z–A'}, then method; unresolved entries last.`;
    });
  }
})();
