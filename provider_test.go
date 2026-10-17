package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvide(t *testing.T) {
	c := NewContainer()

	cat1 := &Cat{Name: "A"}
	dog1 := &Dog{Name: "B"}

	err := c.Provide(cat1)
	assert.Nil(t, err)
	err = c.Provide(dog1)
	assert.Nil(t, err)

	cat2 := &Cat{}
	dog2 := &Dog{}

	err = c.Find(cat2)
	assert.Nil(t, err)
	assert.Equal(t, cat1.Name, cat2.Name)

	err = c.Find(dog2)
	assert.Nil(t, err)
	assert.Equal(t, dog1.Name, dog2.Name)
}

func TestProvideAsInterface(t *testing.T) {
	c := NewContainer()

	cat0 := &Cat{Name: "azure"}
	cat1 := NewPetCat("cat")
	animalCat := NewAnimalCat("ccc")

	err := c.Provide(cat0)
	assert.Nil(t, err)
	err = c.ProvideAs(cat1, (*Pet)(nil))
	assert.Nil(t, err)
	err = c.ProvideAs(animalCat, (*Animal)(nil))
	assert.Nil(t, err)

	cat := &Cat{Name: ""}
	var P Pet
	var A Animal

	err = c.Find(cat)
	assert.Nil(t, err)
	err = c.Find(&P)
	assert.Nil(t, err)
	err = c.Find(&A)
	assert.Nil(t, err)

	assert.Equal(t, cat0.GetName(), cat.Name)
	assert.Equal(t, cat1.GetName(), P.GetName())
	assert.Equal(t, animalCat.GetName(), A.GetName())
}

func TestProvideAsOverride(t *testing.T) {
	c := NewContainer()

	pet1 := NewPetCat("aaa")
	pet2 := NewPetDog("bbb")

	err := c.ProvideAs(pet1, (*Pet)(nil))
	assert.Nil(t, err)
	err = c.ProvideAs(pet2, (*Pet)(nil))
	assert.Nil(t, err)

	var p Pet

	err = c.Find(&p)
	assert.Nil(t, err)

	assert.NotEqual(t, p.GetName(), pet1.GetName())
	assert.Equal(t, p.GetName(), pet2.GetName())
}

func TestProvideInvalidType(t *testing.T) {
	c := NewContainer()

	a := 3

	err := c.Provide(nil)
	assert.NotNil(t, err)

	err = c.Provide(&a)
	assert.NotNil(t, err)

	err = c.Provide(&a)
	assert.NotNil(t, err)

	err = c.ProvideAs(nil, nil)
	assert.NotNil(t, err)

	err = c.ProvideAs(&a, (*Pet)(nil))
	assert.NotNil(t, err)
}

func TestProvideAsNonValue(t *testing.T) {
	c := NewContainer()

	err := c.ProvideAs((*Cat)(nil), (*Pet)(nil))
	assert.NotNil(t, err)

	err = c.ProvideAs((*Cat)(nil), nil)
	assert.NotNil(t, err)
}

func TestProvideAsTargetNonPointInterface(t *testing.T) {
	c := NewContainer()

	err := c.ProvideAs(&Cat{Name: "A"}, Cat{Name: "A"})
	assert.NotNil(t, err)
}

func TestProvideAsWrongImplementation(t *testing.T) {
	c := NewContainer()

	// not implement the interface
	err := c.ProvideAs(&temp{}, (*Pet)(nil))
	assert.NotNil(t, err)
}

// TestProvideFunction 现在验证:Provide/ProvideAs 一律拒绝函数。
// 过去把 func 传进去会被静默注册成"每次解析新建一个"的瞬时工厂,调用方却以为拿到的是
// 共享单例 —— 这是个静默多实例坑。现在收敛为"只接受单例实例(结构体指针)",函数一律报错。
// 瞬时工厂如将来确需,另开显式 API(如 ProvideFactory),不走 Provide。
func TestProvideFunction(t *testing.T) {
	f1 := func() {}
	f2 := func() Cat { return Cat{} }
	f3 := func() (*Cat, error) { return nil, nil }
	f4 := func() *Cat { return &Cat{Name: "A"} }
	f5 := func() Pet { return &Cat{Name: "B"} }

	c := NewContainer()

	// 所有函数形态(无返回/值返回/多返回/*struct 返回/接口返回)一律被 Provide 拒绝。
	for _, f := range []interface{}{f1, f2, f3, f4, f5} {
		err := c.Provide(f)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "functions are not supported")
	}

	// ProvideAs 同样拒绝函数——接口目标 与 *struct 目标都拒(拒绝发生在任何目标类型处理之前)。
	err := c.ProvideAs(f5, (*Pet)(nil))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "functions are not supported")

	err = c.ProvideAs(f4, (*Cat)(nil))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "functions are not supported")

	// 被拒绝后 store 里不应有任何东西:Find 必然 not found。
	cat := Cat{}
	err = c.Find(&cat)
	assert.NotNil(t, err)

	var pet Pet
	err = c.Find(&pet)
	assert.NotNil(t, err)
}

// TestProvideAsStructValueReturnsErrorNotPanic 回归:object 传**结构体值**(而非指针)时,
// ProvideAs 必须返回可读 error,绝不 panic。Cat 用值接收者实现 Pet,故 Cat{} 是个"满足接口的
// 结构体值"——正是"provider 助手误返回 foo{} 而非 &foo{}"的场景(参见 aurora StaticPolicy/
// StaticLimits 曾踩的坑)。修复前:IsNil 先于种类校验执行,对结构体值抛
// "reflect: call of reflect.Value.IsNil on struct Value" 硬 panic(注册即崩、无编译错)。
// 修复后:种类校验先行,返回下面这个 error。本测试能正常跑完(不 panic)本身即验证点。
func TestProvideAsStructValueReturnsErrorNotPanic(t *testing.T) {
	c := NewContainer()

	err := c.ProvideAs(Cat{Name: "value-cat"}, (*Pet)(nil))
	assert.NotNil(t, err, "结构体值应返回 error,不能 panic")
	assert.Contains(t, err.Error(), "object must be pointer of struct")

	// 指针仍正常(未回归)。
	err = c.ProvideAs(&Cat{Name: "ptr-cat"}, (*Pet)(nil))
	assert.Nil(t, err)

	// typed-nil 指针仍被 IsNil 捕获(种类校验通过后 IsNil 依然生效)。
	err = c.ProvideAs((*Cat)(nil), (*Pet)(nil))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "object value cannot be nil")
}

// TestSingletonSameInstanceAcrossResolutions 锁定核心语义:注册一个实例后,多次解析返回的
// 永远是**同一个实例**(单例),不是每次新建/副本。去掉瞬时函数路径后必须保证这点——否则像
// *gorm.DB / redis 这种若被每次解析重建,连接池就不再共享(每请求新连接)。本测试防回归。
func TestSingletonSameInstanceAcrossResolutions(t *testing.T) {
	c := NewContainer()

	cat := &Cat{Name: "singleton"}
	assert.Nil(t, c.Provide(cat))

	var p1, p2 *Cat
	assert.Nil(t, c.Find(&p1))
	assert.Nil(t, c.Find(&p2))
	assert.Same(t, cat, p1) // 解析回来的就是注册进去的那个指针,非副本/新建
	assert.Same(t, p1, p2)  // 多次解析恒为同一实例

	// 接口解析同理:两次解析是同一底层实例。
	pet := NewPetCat("x")
	assert.Nil(t, c.ProvideAs(pet, (*Pet)(nil)))
	var a, b Pet
	assert.Nil(t, c.Find(&a))
	assert.Nil(t, c.Find(&b))
	assert.True(t, a == b)
}
