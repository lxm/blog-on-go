package models

type Relationships struct {
	ObjectId       uint64 `orm:"column(object_id)"`
	TermTaxonomyId uint64 `orm:"column(term_taxonomy_id)"`
	TermOrder      int    `orm:"column(term_order)"`
}
