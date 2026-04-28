---
title: "6章 EKSクラスターにデプロイしてみる"
free: false
---

# はじめに
5章で用意した環境にデプロイしていきます。
具体的には、構成図のECRにdocker imageをpushし、EKSクラスタにmanifestを適用させます。

概念的に表すと以下のようになります。
![](https://storage.googleapis.com/zenn-user-upload/4ee913c95896-20251211.png)


以下の順に行なっていきます。

1. `terrform apply`でAWSリソースを構築する
2. ECRにコンテナイメージをPushする
3. アプリケーションをDeployする

::::message
本章では料金が発生します。
実行は自己責任でお願いいたします。
::::

# `terrform apply`でAWSリソースを構築する

まずは、`terraform plan`を実行して、実行計画を出力しましょう。`./terraform/`にカレントディレクトリを移してから以下を実行してください。
```
terraform plan
```

大量のdiff (全て +)が出力されたかと思います。
最後に、以下のメッセージが出ています。69リソースを「追加」するというPlan結果が得られました。
```sh
Plan: 69 to add, 0 to change, 0 to destroy.
```
新規構築なので期待通りと言えそうです。

では、早速applyしていきましょう。


```
terraform apply
```

実行すると、再度Plan結果のdiffが表示され、実行確認を求められます。
yesで進めましょう。
```sh
Plan: 69 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value:
```

出力結果を貼っておきます。進行中ログは省いています。
```sh
module.eks.aws_cloudwatch_log_group.this[0]: Creating...
module.vpc.aws_vpc.this[0]: Creating...
aws_iam_role.eks_cluster_role: Creating...
aws_iam_role.eks_node_group_role: Creating...
aws_ecr_repository.ecr_repository: Creating...

module.eks.aws_cloudwatch_log_group.this[0]: Creation complete
aws_ecr_repository.ecr_repository: Creation complete
aws_iam_role.eks_node_group_role: Creation complete
aws_iam_role.eks_cluster_role: Creation complete

aws_iam_role_policy_attachment.eks_cluster_AmazonEKSClusterPolicy: Creating...
aws_iam_role_policy_attachment.eks_node_group_AmazonEKSWorkerNodePolicy: Creating...
aws_iam_role_policy_attachment.eks_node_group_ssm: Creating...
aws_iam_role_policy_attachment.eks_node_group_AmazonEKS_CNI_Policy: Creating...
aws_iam_role_policy_attachment.eks_node_group_AmazonEC2ContainerRegistryPullOnly: Creating...

module.vpc.aws_vpc.this[0]: Creation complete
module.vpc.aws_default_route_table.default[0]: Creation complete
module.vpc.aws_internet_gateway.this[0]: Creation complete
module.vpc.aws_subnet.public[0]: Creation complete
module.vpc.aws_subnet.public[1]: Creation complete
module.vpc.aws_subnet.private[0]: Creation complete
module.vpc.aws_subnet.private[1]: Creation complete

module.vpc.aws_route_table.public[0]: Creation complete
module.vpc.aws_route_table.private[0]: Creation complete
module.vpc.aws_route_table_association.public[0]: Creation complete
module.vpc.aws_route_table_association.public[1]: Creation complete
module.vpc.aws_route_table_association.private[0]: Creation complete
module.vpc.aws_route_table_association.private[1]: Creation complete

module.vpc.aws_eip.nat[0]: Creation complete
module.vpc.aws_nat_gateway.this[0]: Creation complete

aws_vpc_endpoint.vpc_endpoint_for_s3: Creation complete
aws_vpc_endpoint.vpc_endpoint_for_ecr_api: Creation complete
aws_vpc_endpoint.vpc_endpoint_for_ecr_dkr: Creation complete

module.eks.module.kms.aws_kms_key.this[0]: Creation complete
module.eks.module.kms.aws_kms_alias.this["cluster"]: Creation complete

module.alb.aws_lb.this[0]: Creation complete
aws_lb_listener.http: Creation complete
aws_lb_target_group.alb_tg_to_ng: Creation complete

module.eks.aws_eks_cluster.this[0]: Creation complete
module.eks.aws_eks_addon.this["vpc-cni"]: Creation complete
module.eks.aws_eks_addon.this["kube-proxy"]: Creation complete

module.eks_node_group.aws_eks_node_group.this[0]: Creation complete
aws_autoscaling_attachment.eks_node_group_attachment: Creation complete

Apply complete! Resources: 69 added, 0 changed, 0 destroyed.
```
Plan通り、69リソースが作成されました。

また、Stateファイルにも変化が見られるはずです。
まさに、今のリモートリソースの状態を反映しています。
```tf:terraform/terraform.tfstate
{
  "version": 4,
  "terraform_version": "1.13.5",
  "serial": 1860,
  "lineage": "ccc76b92-b59f-0be1-2928-9f8cd6cfd64f",
  "outputs": {},
  "resources": [
     // ...
  ],
  "check_results": null
}
```

## マネコンを覗いてみる
1コマンドでリソースが作成されましたと言われても実感が湧かないと思うので、マネコン(AWS マネジメントコンソール)でいくつか覗いてみましょう。

https://aws.amazon.com/jp/console/

### VPC
検索欄でVPCとかで調べてみるとマッチするかと思います。
サブネットに関しても、プライベート・パブリックのが１つずつ、2つのazに配置されているのがわかります。

![](https://static.zenn.studio/user-upload/9e1185790df8-20260426.png)

### EKS
Nodeとして起動させているマシンスペックもt3.microとみえていますね。

![](https://static.zenn.studio/user-upload/39f473a0ec9c-20260426.png)


Podとしてkube-proxyが動いているのが見えたりもします。
![](https://static.zenn.studio/user-upload/84ced70cd698-20260426.png)

アプリケーションをデプロイしたら後でもう一度見てみましょう。


::::message
ルートユーザーでログインした場合、EKSクラスタが存在しているところまでは見えても、以下のようにクラスタの中身を見れないと思います。
これは、アクセスエントリにより権限を指定のIAMユーザーに絞れている証拠と言えます。

![](https://static.zenn.studio/user-upload/15405e79e983-20260426.png)
::::

### ALB

リスナールール / ターゲットグループの関係性もリソースマップで可視化されています。

また、ヘルスチェックに失敗しているのがわかりますね。これは、ヘルスチェックエンドポイントを持つアプリケーションをまだデプロイしていないためです。
こちらも、デプロイ後再度みてみましょう。

![](https://static.zenn.studio/user-upload/3207fdb55279-20260426.png)

----

他も色々みてみてください。


# ECRにコンテナイメージをPushする
Kindクラスターの時と同じく、manifestではデプロイするアプリケーションのコンテナイメージを記述します。
```yml
# k8s/values/production/values.yaml
deployment:
  name: app
  image:
    repository: xxx
    tag:  xxxx

# manifestには実際以下のように展開される
# k8s/charts/http-server/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.deployment.image.repository }}:{{ .Values.deployment.image.tag }}"
```

EKSでデプロイを行うために、kubeletがPullするためのコンテナイメージをECRに用意していきましょう。

手順としては、手元でコンテナイメージをビルドして、ECRにPushするという流れになります。

まずは、以下のコマンドでどのIAMユーザーとしてAWSにログインしているかを確認するところから始めましょう。
```
aws sts get-caller-identity
```

Push先のECRのURIを控えましょう。
```sh
ECR_REPO_URI=$(aws ecr describe-repositories --repository-names eks-sandbox-ecr-repository --query 'repositories[0].repositoryUri' --output text)
```
```sh
#　確認
echo $ECR_REPO_URI 
846869429016.dkr.ecr.ap-northeast-1.amazonaws.com/eks-sandbox-ecr-repository
```

次に、対象のECRにログインします。
```sh
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin $ECR_REPO_URI
```
`Login Succeeded`と出ればOKです。

Pushするための準備は整いました。
Kindの時と同じようにコンテナをビルドしましょう。

```sh
TIMESTAMP=$(date +"%Y%m%d%H%M%S")
docker build -t app:$TIMESTAMP .
```

ビルドし終わったら、ECRへのPushのためにイメージ名をもう一つつけます。
```sh
docker tag app:$TIMESTAMP $ECR_REPO_URI:$TIMESTAMP
```

イメージを確認しておきます。
```sh
% docker image ls
REPOSITORY                                                                     TAG              IMAGE ID       CREATED         SIZE
846869429016.dkr.ecr.ap-northeast-1.amazonaws.com/eks-sandbox-ecr-repository   20260426020650   4bc4d9dd37dc   3 minutes ago   109MB
```

このイメージは `846869429016.dkr.ecr.ap-northeast-1.amazonaws.com` ECRの
`eks-sandbox-ecr-repository` というリポジトリに対応する名前が付けられており、
この名前を使って `docker push` を実行すると、そのリポジトリにイメージがpushされます。

では、これをPushしましょう。
```sh
docker push $ECR_REPO_URI:$TIMESTAMP
```
出力として、以下のように出ていると思います。期待通りですね。
```sh
The push refers to repository [846869429016.dkr.ecr.ap-northeast-1.amazonaws.com/eks-sandbox-ecr-repository]
```

マネコンからも確認しておきましょう。
![](https://static.zenn.studio/user-upload/2528e4b2b5f6-20260426.png)

また、production用のvalues.yamlにECRリポジトリ名 + イメージタグを記載しておきましょう。
```yaml:k8s/values/production/values.yaml
deployment:
  name: app
  image:
    repository: 846869429016.dkr.ecr.ap-northeast-1.amazonaws.com/eks-sandbox-ecr-repository
    tag: "20260426020650"
```

# アプリケーションをDeployする
まずは、`kubectl`でEKSクラスタ（eks-sandbox）に接続できるように以下のコマンドを実行します。
```sh
aws eks update-kubeconfig --region ap-northeast-1 --name eks-sandbox --alias eks-sandbox
```
これにより、EKSクラスタの接続先情報（APIサーバーのエンドポイントや証明書）と、IAMを用いた認証設定がローカルのkubeconfig（`~/.kube/config`）に追加されます。

https://docs.aws.amazon.com/ja_jp/eks/latest/userguide/create-kubeconfig.html

これにより、`eks-sandbox`というコンテキスト（どのユーザーでどのクラスタに接続するか）が作成されます。
```sh
% kubectl config get-contexts
CURRENT   NAME                                                          CLUSTER                                                       AUTHINFO                                                      NAMESPACE
*         eks-sandbox                                                   arn:aws:eks:ap-northeast-1:846869429016:cluster/eks-sandbox   arn:aws:eks:ap-northeast-1:846869429016:cluster/eks-sandbox   
          kind-sandbox                                                  kind-sandbox                                                  kind-sandbox                               
```

helmfileには以下のように設定していたので、`-e production`とすれば、`eks-sandbox`コンテキストでデプロイされるようになっています。
```yaml:k8s/helmfile.yaml.gotmpl
environments:
  development:
    kubeContext: kind-sandbox

  production:
    kubeContext: eks-sandbox 

---
releases:
  - name: app
    namespace: "app-{{ .Environment.Name }}"
    chart: ./charts/http-server
    values:
      - values/{{ .Environment.Name }}/values.yaml
    installed: true
    atomic: false # Rollback on failure is disabled for learning purposes
```

それでは、`/k8s`にカレントディレクトリを移してデプロイを実行します。
```sh
helmfile -e production apply
```

# デプロイができているか確認する

正常にPodが作成されているのか、`kubectl`で確認してみます。
`app-http-server-85b8bbd4c7-fsw97`PodがRUNNINGなので良さそうですね。
```sh
% kubectl get pods -A
NAMESPACE        NAME                               READY   STATUS    RESTARTS   AGE
app-production   app-http-server-85b8bbd4c7-fsw97   2/2     Running   0          34s
kube-system      aws-node-6ttjb                     2/2     Running   0          89m
kube-system      kube-proxy-q427w                   1/1     Running   0          89m
kube-system      metrics-server-59b569559b-4lsxj    1/1     Running   0          34s
```

::::message
`kube-system`namespaceのpodがKindで見た時よりもPod数が少ないと思いませんか？

Podが乗っているNodeも出力されるようにすると、以下のようになりました。
```sh
% kubectl get pods -A -o custom-columns=NAMESPACE:.metadata.namespace,NAME:.metadata.name,NODE:.spec.nodeName --sort-by='.spec.nodeName'

NAMESPACE        NAME                               NODE
app-production   app-http-server-85b8bbd4c7-fsw97   ip-10-0-3-216.ap-northeast-1.compute.internal
kube-system      aws-node-6ttjb                     ip-10-0-3-216.ap-northeast-1.compute.internal
kube-system      kube-proxy-q427w                   ip-10-0-3-216.ap-northeast-1.compute.internal
kube-system      metrics-server-59b569559b-4lsxj    ip-10-0-3-216.ap-northeast-1.compute.internal
```

`ip-***.ap-northeast-1.compute.internal`はWorker Nodeですね。Controle PlaneのPodが一切表示されていません。

これは、Controle PlaneがAWSマネージドであることによるものだと筆者は理解しています。
::::

また、マネコンでも見てみると、ちゃんとこちらでも見えますね。イベントも追えます。
![](https://static.zenn.studio/user-upload/1f5a32d48be3-20260426.png)

Pod内コンテナも参照できます。利用イメージも、http-serverの方は先ほどPushしたものになっていますね。
![](https://static.zenn.studio/user-upload/55362f6d848c-20260426.png)

また、デプロイ前は失敗していたALBのヘルスチェックも、成功しています。
![](https://static.zenn.studio/user-upload/8348f3bc232a-20260426.png)

# アプリケーションサーバーにアクセスしてみる
さて、最後にデプロイされているTODOアプリケーションにアクセスしてみましょう。

ALBが公開しているDNS名を取得します。
```sh
ALB_DNS=$(aws elbv2 describe-load-balancers --query 'LoadBalancers[?LoadBalancerName==`eks-sandbox-alb`].DNSName' --output text)
```

筆者の場合、`eks-sandbox-alb-524784966.ap-northeast-1.elb.amazonaws.com`のようです。
```sh
# 確認
echo $ALB_DNS
eks-sandbox-alb-524784966.ap-northeast-1.elb.amazonaws.com
```

ブラウザのアドレスバーに、`$ALB_DNS`の値をドメインとして`http://$ALB_DNS/todos`の形式で打ち込んでみると、以下のようにTODOアプリにアクセスできました。
![](https://static.zenn.studio/user-upload/deff89bbdb63-20260426.png)

ログも正常に流れていることを確認しました。
```sh
% kubectl logs app-http-server-85b8bbd4c7-fsw97 -n app-production | grep todos

2026/04/25 17:58:52 [::1]:56670 GET /todos 200
2026/04/25 18:00:32 [::1]:33240 POST /todos 303
2026/04/25 18:00:32 127.0.0.1:39564 GET /todos 200
```

# 環境のお掃除
色々触って確認できたら、最後にお掃除をしましょう。これを忘れると料金が発生し続けるのでご注意ください。

まずはアプリケーションを削除します。
`./k8s`で行ってください。
```sh
helmfile -e production destroy
```

そして、AWSリソース(インフラ)を削除します。
`./terraform`で行ってください。
```sh
terraform destroy
```