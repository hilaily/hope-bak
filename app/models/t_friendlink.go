package models

var FriendlinkModel *Friendlink

type Friendlink struct {
	Id   int64  `xorm:"pk autoincr INT(11)"`
	Name string `xorm:"not null VARCHAR(100)"`
	Link string `xorm:"not null VARCHAR(255)"`
}

//func (f *Friendlink) List() []*Friendlink {
//	f := &Friendlink{}
//	support.XORM.WHERE()
//}
