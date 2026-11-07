function calendar(serviceId) {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return {
    month: new Date(today.getFullYear(), today.getMonth(), 1),
    selected: null,
    get title() {
      return this.month.toLocaleDateString('ru-RU', { month: 'long', year: 'numeric' })
    },
    get days() {
      const y = this.month.getFullYear(), m = this.month.getMonth()
      const offset = (this.month.getDay() + 6) % 7
      const count = new Date(y, m + 1, 0).getDate()
      return [...Array(offset).fill(null), ...Array.from({ length: count }, (_, i) => new Date(y, m, i + 1))]
    },
    shift(n) {
      this.month = new Date(this.month.getFullYear(), this.month.getMonth() + n, 1)
    },
    past(d) {
      return d < today
    },
    iso(d) {
      return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
    },
    pick(d) {
      this.selected = this.iso(d)
      htmx.ajax('GET', `/booking/slots?service_id=${serviceId}&date=${this.selected}`, '#slots')
    },
  }
}
