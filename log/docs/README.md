# Terraform Plugin Framework
https://www.terraform.io/docs/plugin/framework/

- Terraform Plugin SDKv2に代わる新しいSDKらしい

## note
### Provider
- Terraform CoreとgRPCで通信するサーバ
- Terraform CLIでダウンロードする
- `GetSchema`, `Configure`, `GetResources`, `GetDataSources`の4メソッドがinterfaceにて定義されている

GetSchema
- APIによって認証するために利用されるクレデンシャルやエンドポイント等を集めるのに使うメソッド
- immutable

Configure
- APIクライアントの作成、Provider interfaceを実装する型へのデータの保存に使うメソッド
- Terraform(Core)がユーザがProviderのConfigurationブロックで指定した値をProviderに対して送った時に、Providerのライフサイクルの最初で呼ばれる
- 実行時に値が確定せずUnknownが入ることもあるので、ハンドリングする必要があることに注意

GetResources
- resource typesのmapを返却するメソッド
- immutable

GetDataSources
- data source typesのmapを返却するメソッド
- immutable

### Resources
- compute instanceやアクセスポリシ、ディスク的なオブジェクトを指す
- `GetSchema`と`NewResource`の2メソッドがinterfaceにて定義されている

GetSchema
- 略

NewResource
- [Resource型](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#Resource)のファクトリメソッド
- ProviderのConfigureメソッドが実行された後に実行される

実装例
```go
// ResourceType 
type computeInstanceResourceType struct{}

func (c computeInstanceResourceType) GetSchema(_ context.Context) (tfsdk.Schema,
    diag.Diagnostics) {
    return tfsdk.Schema{
        Attributes: map[string]tfsdk.Attribute{
            "name": {
                Type: types.StringType,
                Required: true,
            },
        },
    }, nil
}

func (c computeInstanceResourceType) NewResource(_ context.Context,
    p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
    return computeInstanceResource{
        client: p.(*provider).client,
    }, nil
}
```

#### Resourceの定義
- ResouceはResource型の**単一の**インスタンスを対象としている

Create
- 必要なAPI呼び出しを作成する
  - Resourceを作るため
  - TerraformのstateにResourceのデータを永続化するため

1. `tfsdk.CreateResourceRequest`からplanされたデータを読む
2. Resource typeの`NewResource`メソッドによってResourceに設定されたAPIクライアントを使う
3. Stateを書き込む


Read
- APIのResourceの最新の状態を反映するために、Terraformの状態を更新する

Update
- 必要なAPI呼び出しを作成する
  - 構成にマッチするようなResourceを編集するため
  - TerrafromのstateにResourceのデータを永続化するため
- Createと同様に処理
  - `tfsdk.UpdateResourceRequest`

Delete
- 必要なAPI呼び出しを作成する
  - Resourceを削除するため
  - その後、TerraformのstateからそのResourceを除去する
- Createと同様に処理

ImportState
- 初期のTerrafromのstateを作成する
  - `terraform import`コマンドを介して管理されているリソースを持ってくるため
- Createと同様に処理

#### 追加するとき
- mapのkeyをproviderのprefixを含めたリソースの名前にするべき
- mapのvalueをresource typeのインスタンスにするべき

### Data Sources
- 外部データをTerraformに参照させるための抽象化したもの
- resourceとの違いは、Terraformがこのリソースを管理しないこと
- 多分使わない

### Schema
- Terraform Configuration blockの定数を定義するもの
  - 何のフィールドをprovider, resource, data sourceが持っているかを定義する
  - それらのフィールドに関するTerraformのメタデータを与える
- 単一のvalue, list, mapが定義できる
- フィールドのrequired, optional, computed(providerにcomputeされるか: boolean), sensitive(機密情報か否か: boolean; local stateではplain-textのJSON, remote stateではその環境に依存する)

### Types
- Stateを持つ
- 非同期ではないPromise的な感じ?

---

- Null
- Unknown
  - Terraformでは実行順序が存在する
  - あるprovider, resource, data sourceが他のproviderの値に依存する場合、値は解決されていない
  - その時に使われる型
- String(string)
- Int64(int64)
- Float64(float64)
- Number(big.Float)
- Bool(bool)
- List(types.List)
  - 異なる型のvalueを格納できる順序つきcollection
- Map(types.Map)
  - uniqueなstringのindexによって異なる型のvalueを格納できる非順序つきcollection
- Object
  - uniqueで事前に指定された他の型のvalueを格納できる非順序つきcollection
- Set
  - データ構造のsetと同じでuniqueな値のみが保持される非順序つきcollection
- 自前で定義する型
- attr.Type, AttributePathStepper, attr.Value
  - インタフェースを実装していれば自前で定義できる
  - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/attr#Type

### エラーハンドリング
既に存在しているdiagnosticsに対してerrorを追加するとき
- Append
  - t.Error()的な
- HasError
  - t.Fatal()的な

---

diagnosticを新しく作るとき
- AddError(summary, detail string)
- AddWarning(summary detail string)

---

AttributePathも指定したい場合はAddAttribute hogehogeを使う

