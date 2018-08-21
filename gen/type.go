package gen

import (
	"strings"

	"github.com/dave/jennifer/jen"
)

func TypeFileContent(typeName string) *jen.File {
	f := jen.NewFile(packageName(typeName))

	f.Add(Defs(typeName))
	f.Add(FromT(typeName))

	f.Func().Id("Collect").Params(jen.Id("iter").Id("Iter")).Params(jen.Index().Id(typeName), jen.Error()).Block(
		jen.Return(jen.Id("collect").Call(jen.Id("iter"))),
	)

	f.Func().Params(jen.Id("f").Op("*").Id("Functor")).Id("Collect").Params().Params(jen.Index().Id(typeName), jen.Error()).Block(
		jen.Return(jen.Id("collect").Call(jen.Id("f").Dot("iter"))),
	)

	f.Func().Id("Collapse").Params(jen.Id("iter").Id("Iter")).Index().Id(typeName).Block(
		jen.Return(jen.Id("collapse").Call(jen.Id("iter"))),
	)

	f.Func().Params(jen.Id("f").Op("*").Id("Functor")).Id("Collapse").Params().Index().Id(typeName).Block(
		jen.Return(jen.Id("collapse").Call(jen.Id("f").Dot("iter"))),
	)

	f.Func().Id("Fold").Params(jen.Id("iter").Id("Iter"), jen.Id("initial").Id(typeName), jen.Id("op").Id("foldErrFunc")).Params(jen.Id(typeName), jen.Error()).Block(
		jen.List(jen.Id("result"), jen.Id("err")).Op(":=").Id("fold").Call(jen.Id("iter"), jen.Id("T").Call(jen.Id("initial")), jen.Id("op")),
		jen.Return(jen.Id("fromT").Call(jen.Id("result")), jen.Id("err")),
	)

	f.Func().Params(jen.Id("f").Op("*").Id("Functor")).Id("Fold").Params(jen.Id("initial").Id(typeName), jen.Id("op").Id("foldErrFunc")).Params(jen.Id(typeName), jen.Error()).Block(
		jen.Return(jen.Id("Fold").Call(jen.Id("f").Dot("iter"), jen.Id("initial"), jen.Id("op"))),
	)

	f.Func().Id("Roll").Params(jen.Id("iter").Id("Iter"), jen.Id("initial").Id(typeName), jen.Id("op").Id("foldFunc")).Id(typeName).Block(
		jen.Return(jen.Id("fromT").Call(jen.Id("roll").Call(jen.Id("iter"), jen.Id("T").Call(jen.Id("initial")), jen.Id("op")))),
	)

	f.Func().Params(jen.Id("f").Op("*").Id("Functor")).Id("Roll").Params(jen.Id("initial").Id(typeName), jen.Id("op").Id("foldFunc")).Id(typeName).Block(
		jen.Return(jen.Id("Roll").Call(jen.Id("f").Dot("iter"), jen.Id("initial"), jen.Id("op"))),
	)

	f.Func().Id("Transmute").Params(jen.Id("v").Interface()).Id(typeName).Block(
		jen.List(jen.Id("result"), jen.Id("ok")).Op(":=").Id("v").Assert(jen.Id(typeName)),
		jen.If(jen.Op("!").Id("ok")).Block(
			jen.Panic(jen.Qual("fmt", "Sprintf").Call(jen.Lit("could not transmute: %v"), jen.Id("v"))),
		),
		jen.Return(jen.Id("result")),
	)

	f.Func().Id("asMapErrFunc").Params(jen.Id("f").Id("mapFunc")).Id("mapErrFunc").Block(
		jen.Return(jen.Func().Params(jen.Id("v").Id(typeName)).Params(jen.Id(typeName), jen.Error()).Block(
			jen.Return(jen.Id("f").Call(jen.Id("v")), jen.Nil()),
		)),
	)

	f.Func().Id("asFilterErrFunc").Params(jen.Id("f").Id("filterFunc")).Id("filterErrFunc").Block(
		jen.Return(jen.Func().Params(jen.Id("v").Id(typeName)).Params(jen.Bool(), jen.Error())).Block(
			jen.Return(jen.Id("f").Call(jen.Id("v")), jen.Nil()),
		),
	)

	f.Func().Id("asFoldErrFunc").Params(jen.Id("f").Id("foldFunc")).Id("foldErrFunc").Block(
		jen.Return(jen.Func().Params(jen.Id("v").Id(typeName), jen.Id("w").Id(typeName)).Params(jen.Id(typeName), jen.Error())).Block(
			jen.Return(jen.Id("f").Call(jen.Id("v"), jen.Id("w")), jen.Nil()),
		),
	)

	return f
}

func Defs(typeName string) *jen.Statement {
	return jen.Type().Defs(
		jen.Id("T").Id(typeName),
		jen.Id("tSlice").Index().Id(typeName),
		jen.Id("mapFunc").Func().Params(jen.Id(typeName)).Id(typeName),
		jen.Id("mapErrFunc").Func().Params(jen.Id(typeName)).Params(jen.Id(typeName), jen.Error()),
		jen.Id("foldFunc").Func().Params(jen.Id(typeName), jen.Id(typeName)).Id(typeName),
		jen.Id("foldErrFunc").Func().Params(jen.Id(typeName), jen.Id(typeName)).Params(jen.Id(typeName), jen.Error()),
		jen.Id("filterFunc").Func().Params(jen.Id(typeName)).Bool(),
		jen.Id("filterErrFunc").Func().Params(jen.Id(typeName)).Params(jen.Bool(), jen.Error()),
		jen.Id("transformFunc").Func().Params(jen.Interface()).Params(jen.Id(typeName), jen.Error()),
	)
}

func FromT(typeName string) *jen.Statement {
	body := jen.Return(jen.Id(typeName).Call(jen.Id("t")))

	if strings.HasPrefix(typeName, "*") {
		body = jen.Return(jen.Id("t"))
	}

	return jen.Func().Id("fromT").Params(jen.Id("t").Id("T")).Id(typeName).Block(body)
}

func packageName(typeName string) string {
	if typeName == "interface{}" {
		return "finterface"
	}
	if strings.HasPrefix(typeName, "*") {
		return "fp" + typeName[1:len(typeName)]
	}
	return "f" + typeName
}
