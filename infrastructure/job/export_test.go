package job

type DataSource = dataSource

func (d *DataSource) SetDB(db LevelDB) (cancel func()) {
	tmp := d.db
	d.db = db
	return func() {
		d.db = tmp
	}
}
