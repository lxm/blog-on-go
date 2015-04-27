package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Posts struct {
	Id                  int       `orm:"column(ID);auto"`
	PostAuthor          uint64    `orm:"column(post_author)"`
	PostDate            time.Time `orm:"column(post_date);type(datetime)"`
	PostDateGmt         time.Time `orm:"column(post_date_gmt);type(datetime)"`
	PostContent         string    `orm:"column(post_content)"`
	PostTitle           string    `orm:"column(post_title)"`
	PostExcerpt         string    `orm:"column(post_excerpt)"`
	PostStatus          string    `orm:"column(post_status);size(20)"`
	CommentStatus       string    `orm:"column(comment_status);size(20)"`
	PingStatus          string    `orm:"column(ping_status);size(20)"`
	PostPassword        string    `orm:"column(post_password);size(20)"`
	PostName            string    `orm:"column(post_name);size(200)"`
	ToPing              string    `orm:"column(to_ping)"`
	Pinged              string    `orm:"column(pinged)"`
	PostModified        time.Time `orm:"column(post_modified);type(datetime)"`
	PostModifiedGmt     time.Time `orm:"column(post_modified_gmt);type(datetime)"`
	PostContentFiltered string    `orm:"column(post_content_filtered)"`
	PostParent          uint64    `orm:"column(post_parent)"`
	Guid                string    `orm:"column(guid);size(255)"`
	MenuOrder           int       `orm:"column(menu_order)"`
	PostType            string    `orm:"column(post_type);size(20)"`
	PostMimeType        string    `orm:"column(post_mime_type);size(100)"`
	CommentCount        int64     `orm:"column(comment_count)"`
}

func (t *Posts) TableName() string {
	return "posts"
}

func init() {
	orm.RegisterModel(new(Posts))
}

// AddPosts insert a new Posts into database and returns
// last inserted Id on success.
func AddPosts(m *Posts) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPostsById retrieves Posts by Id. Returns error if
// Id doesn't exist
func GetPostsById(id int) (v *Posts, err error) {
	o := orm.NewOrm()
	v = &Posts{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPosts retrieves all Posts matches certain condition. Returns empty list if
// no records exist
func GetAllPosts(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Posts))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Posts
	qs = qs.OrderBy(sortFields...)
	if _, err := qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdatePosts updates Posts by Id and returns error if
// the record to be updated doesn't exist
func UpdatePostsById(m *Posts) (err error) {
	o := orm.NewOrm()
	v := Posts{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePosts deletes Posts by Id and returns error if
// the record to be deleted doesn't exist
func DeletePosts(id int) (err error) {
	o := orm.NewOrm()
	v := Posts{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Posts{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
