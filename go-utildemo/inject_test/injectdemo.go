package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/facebookgo/inject"
)

// 在HomePlanetRenderApp上构造两个API类
type HomePlanetRenderApp struct {

	//下面的标签向注入库指示这些字段
	//表示有资格注射。 他们没有指定任何选项，并将会
	//导致为每个API创建的单例实例。

	NameAPI   *NameAPI   `inject:""`
	PlanetAPI *PlanetAPI `inject:""`
}

func (a *HomePlanetRenderApp) Render(id uint64) string {
	return fmt.Sprintf(
		"%s is from the planet %s.",
		a.NameAPI.Name(id),
		a.PlanetAPI.Planet(id),
	)
}

// 伪造NameAPI
type NameAPI struct {
	//在PlanetAPI中下面和这里下面我们将标签添加到接口值中。该值不能自动创建（按定义），因此必须显式提供给图形。
	HTTPTransport http.RoundTripper `inject:""`
}

func (n *NameAPI) Name(id uint64) string {
	return "Spock"
}

// Our fake Planet API.
type PlanetAPI struct {
	HTTPTransport http.RoundTripper `inject:""`
}

func (p *PlanetAPI) Planet(id uint64) string {
	return "Vulcan"
}

func main() {
	// 通常，应用程序将只有一个对象图形，您将创建它并在主要功能中使用它：
	var g inject.Graph

	//我们向我们提供了两个“种子”对象，一个是我们的空的
	//我们自动填充的HomePlanetRenderApp实例，并且需要满足PlanetAPI对HTTPTransport依赖性。http.RoundTripper是接口所以我们只需要提供一个派生类即可：

	var a HomePlanetRenderApp
	err := g.Provide(
		&inject.Object{Value: &a},
		&inject.Object{Value: http.DefaultTransport},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	//这里Populate调用正在创建NameAPI＆
	// PlanetAPI，并将HTTPTransport设置为
	// http.DefaultTransport提供如下：

	if err := g.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	//有一个简化的case的简写API，它结合了
	//以上三个调用可用作inject.Populate：
	//inject.Populate（＆a，http.DefaultTransport）
	//上述API显示了还允许使用的底层API更复杂场景的命名实例。
	fmt.Println(a.Render(42))
}
