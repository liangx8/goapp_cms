package entity

type (
	User struct {
		Seq      uint64 `name:"seq" primary:""`
		Name     string `name:"name"`
		Pwd      []byte `name:"pwd"`
		Password string
		Active   bool   `name:"active"`
		Remark   string `name:"remark"`
		Updated  uint64 `name:"updated"`
	}
	// 参与者，可以是供应商，客户，分销商，以及最终消费者
	Participate struct {
		Seq    uint64 `name:"seq" primary:""`
		Name   string `name:"name"`
		Remark string `name:"remark"`
		Role   string `name:"role"`
	}
	// 配件，产品，原料
	Item struct {
		Seq   uint64 `name:"seq" primary:""`
		Icode string `name:"icode" desc:"条码值,item的编号等"`
		Name  string `name:"name"`
	}
)
