package dist
import(
	"cloud.google.com/go/storage"

	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
)
func allEntity(ctx context.Context) ([]Entity,error){
	var ary []Entity
	cli,err := storage.NewClient(ctx)
	if err != nil {
		return nil,err
	}
	defer cli.Close()

	bucket := cli.Bucket(BUCKETNAME)
	oh := bucket.Object(TOUCH_DB)
	objr,err := oh.NewReader(ctx)
	if err != nil {
		return nil,err
	}
	defer objr.Close()
	de:=yaml.NewDecoder(objr)
	if err=de.Decode(&ary);err != nil {
		return nil,err
	}
	
	return ary,nil
}
func add(ctx context.Context,e Entity) error{
	ens,err:=allEntity(ctx)
	if err != nil {
		ens = make([]Entity,1)
		ens[0]=e
	} else {
		ens = append(ens,e)
	}
	
	return save(ctx, ens)
}
func save(ctx context.Context, ens []Entity) error{
	cli,err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer cli.Close()

	bucket := cli.Bucket(BUCKETNAME)
	oh := bucket.Object(TOUCH_DB)
	objw:=oh.NewWriter(ctx)
	defer objw.Close()
	yamlEn:=yaml.NewEncoder(objw)
	defer yamlEn.Close()
	err = yamlEn.Encode(ens)
	if err != nil { return err}
	return nil
}
const (
	BUCKETNAME="pfa.rc-greed.com"
	TOUCH_DB="touch-db.yaml"
)
