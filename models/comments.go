package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Comments struct {
	Id                 int       `orm:"column(comment_ID);auto"`
	CommentPostID      uint64    `orm:"column(comment_post_ID)"`
	CommentAuthor      string    `orm:"column(comment_author)"`
	CommentAuthorEmail string    `orm:"column(comment_author_email);size(100)"`
	CommentAuthorUrl   string    `orm:"column(comment_author_url);size(200)"`
	CommentAuthorIP    string    `orm:"column(comment_author_IP);size(100)"`
	CommentDate        time.Time `orm:"column(comment_date);type(datetime)"`
	CommentDateGmt     time.Time `orm:"column(comment_date_gmt);type(datetime)"`
	CommentContent     string    `orm:"column(comment_content)"`
	CommentKarma       int       `orm:"column(comment_karma)"`
	CommentApproved    string    `orm:"column(comment_approved);size(20)"`
	CommentAgent       string    `orm:"column(comment_agent);size(255)"`
	CommentType        string    `orm:"column(comment_type);size(20)"`
	CommentParent      uint64    `orm:"column(comment_parent)"`
	UserId             uint64    `orm:"column(user_id)"`
}

func (t *Comments) TableName() string {
	return "comments"
}

func init() {
	orm.RegisterModel(new(Comments))
}

// AddComments insert a new Comments into database and returns
// last inserted Id on success.
func AddComments(m *Comments) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCommentsById retrieves Comments by Id. Returns error if
// Id doesn't exist
func GetCommentsById(id int) (v *Comments, err error) {
	o := orm.NewOrm()
	v = &Comments{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllComments retrieves all Comments matches certain condition. Returns empty list if
// no records exist
func GetAllComments(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Comments))
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

	var l []Comments
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

// UpdateComments updates Comments by Id and returns error if
// the record to be updated doesn't exist
func UpdateCommentsById(m *Comments) (err error) {
	o := orm.NewOrm()
	v := Comments{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteComments deletes Comments by Id and returns error if
// the record to be deleted doesn't exist
func DeleteComments(id int) (err error) {
	o := orm.NewOrm()
	v := Comments{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Comments{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
