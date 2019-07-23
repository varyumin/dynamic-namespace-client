# Kubectl-dynamicns

![kubectl logo](https://raw.githubusercontent.com/kubernetes/kubectl/master/images/kubectl-logo-medium.png)

`kubectl-dynamic` - это клиентская часть. Которая упрощает взаимодействие с `Dynamic Namespace Server`. Взаимодействие осуществляется по HTTP запросами к REST API сервера. В разработке использовался фреймворк [cobra](https://github.com/spf13/cobra)

## Dependencies

Зависимости управляются с помощью [dep](https://github.com/golang/dep). При необходимости обратитесь к документации.

## Build
`make build` - соберет Вам бинарный файл в текущей директории под именем `kubectl-dynamicns`

`make docker-build` - Собирает `kubectl-dynamicns` и кладет его в docker.

## Getting started

### Running kubectl-dynamicns
#### HELP. kubectl-dynamicns --help
`kubectl-dynamicns` - Имеет 3 глобальных, **ОБЯЗАТЕЛЬНЫХ** флагов `--cluster`, `--password` и `--user`. Их необходимо передавать во **ВСЕХ** запросах, либо их можно указать через `Environment OS`.  
```
Usage:
  kubectl-dynamicns [command]

Available Commands:
  create      Create dynamic namespace
  delete      Delete dynamic namespace
  help        Help about any command
  list        List all your dynamic namespaces

Flags:
  -c, --cluster string    Cluster ENV. Environment CLUSTER (default environment "CLUSTER")
  -P, --password string   Password from LDAP. Environment LDAP_PASSWORD (default environment "LDAP_PASSWORD")
  -U, --user string       User from LDAP. Environment LDAP_USER (default environment "LDAP_USER")
```

#### HELP. kubectl-dynamicns create --help
```
Create dynamic namespace

Usage:
  kubectl-dynamicns create [flags]

Flags:
  -g, --groups stringArray    Access groups
  -h, --help                  help for create
      --limit.cpu float32     Limnit CPU
      --limit.memory int      Limnit memory
      --request.cpu float32   Request CPU
      --request.memory int    Request memory
  -d, --time-delete string    Timespamp delete
  -u, --users stringArray     Access users

Global Flags:
  -c, --cluster string    Cluster ENV. Environment CLUSTER (default environment "CLUSTER")
  -P, --password string   Password from LDAP. Environment LDAP_PASSWORD (default environment "LDAP_PASSWORD")
  -U, --user string       User from LDAP. Environment LDAP_USER (default environment "LDAP_USER")

```

`kubectl-dynamicns create` - Создает `dynamic-namespace`. 

#### Обязательные флаги:
* `--limit.cpu` - **float32**.Сколько `limits.cpu` будет добавлено в [Resource Quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/#compute-resource-quota).
* `--limit.memory` - **int32** и измеряется в `Mb` .Сколько `limits.memory` будет добавлено в [Resource Quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/#compute-resource-quota). 

#### Не обязательные флаги:
* `--time-delete`  - **string** в формате [RFC3339](https://www.ietf.org/rfc/rfc3339.txt). Дата когда удалить `dynamic namespace`. **Если не указать дату `dynamic namespace` будет создан с TTL +1 день**
* `--users` - **stringArray**. Указать список пользователей из LDAP кому необходимо выдать права `admins` на  `dynamic namespace`. **По умолчанию добавляется  владелец `dynamic namespace` он же `--user string`**
* `--groups` - **stringArray**. Указать список групп из LDAP кому необходимо выдать права `admins` на  `dynamic namespace`. **По умолчанию добавляется  группа `devops`**
* `--request.cpu` - **float32**.Сколько `requests.cpu` будет добавлено в [Resource Quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/#compute-resource-quota). **По умолчанию `requests.cpu` = `limits.cpu`**
* `--request.memory` - **int32** и измеряется в `Mb` .Сколько `requests.memory` будет добавлено в [Resource Quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/#compute-resource-quota). **По умолчанию `requests.memory` = `limits.memory`**

#### Example:
##### Short version:
```
export CLUSTER=dev2
export LDAP_USER=v.a.ryumin
export LDAP_PASSWORD=P@ssword
./kubectl-dynamicns create --limit.cpu=1 --limit.memory=1000

{
  "name": "dynamic-namespace-rfj72",
  "date_create": "2019-07-22T11:43:42Z",
  "date_delete": "2019-07-23T11:43:41Z",
  "owner": "v.a.ryumin",
  "resource_quota_name": "default",
  "limit": {
    "cpu": "1",
    "memory": "1000Mi"
  },
  "request": {
    "cpu": "1",
    "memory": "1000Mi"
  },
  "role_binding_name": "admins",
  "access": {
    "user": ["v.a.ryumin"],
    "group": ["devops"]
  }
}

```
##### Full version:
```
./kubectl-dynamicns --cluster=dev2 --user=v.a.ryumin --password=P@ssword \
    create --time-delete=2019-07-27T11:00:00Z --users=a.invanov,b.petrov --groups=team,cicd --limit.cpu=2 --limit.memory=2000 --request.cpu=1 --request.memory=1000 

{
  "name": "dynamic-namespace-zpz52",
  "date_create": "2019-07-22T11:49:13Z",
  "date_delete": "2019-07-27T11:00:00Z",
  "owner": "v.a.ryumin",
  "resource_quota_name": "default",
  "limit": {
    "cpu": "2",
    "memory": "2000Mi"
  },
  "request": {
    "cpu": "1",
    "memory": "1000Mi"
  },
  "role_binding_name": "admins",
  "access": {
    "user": ["a.invanov,b.petrov", "v.a.ryumin"],
    "group": ["devops", "team", "cicd"]
  }
}

```

#### HELP. kubectl-dynamicns list --help
`kubectl-dynamicns list` - получить таблицу  с `dynamic namespace`, которые были созданы тобой
```
List all your dynamic namespaces

Usage:
  kubectl-dynamicns list [flags]

Flags:
  -h, --help   help for list

Global Flags:
   -c, --cluster string    Cluster ENV. Environment CLUSTER (default environment "CLUSTER")
   -P, --password string   Password from LDAP. Environment LDAP_PASSWORD (default environment "LDAP_PASSWORD")
   -U, --user string       User from LDAP. Environment LDAP_USER (default environment "LDAP_USER")

```
#### Example:
##### Short version:
```
export CLUSTER=dev2
export LDAP_USER=v.a.ryumin
export LDAP_PASSWORD=P@ssword
./kubectl-dynamicns list

NAMESPACE                       OWNER                           CREATE                          DELETE                          
dynamic-namespace-rfj72         v.a.ryumin                      2019-07-22 11:43:42 +0000 UTC   2019-07-23 11:43:41 +0000 UTC   
dynamic-namespace-zpz52         v.a.ryumin                      2019-07-22 11:49:13 +0000 UTC   2019-07-27 11:00:00 +0000 UTC  
```
##### Full version:
```
./kubectl-dynamicns --cluster=dev2 --user=v.a.ryumin --password=P@ssword list

NAMESPACE                       OWNER                           CREATE                          DELETE                          
dynamic-namespace-rfj72         v.a.ryumin                      2019-07-22 11:43:42 +0000 UTC   2019-07-23 11:43:41 +0000 UTC   
dynamic-namespace-zpz52         v.a.ryumin                      2019-07-22 11:49:13 +0000 UTC   2019-07-27 11:00:00 +0000 UTC  
```

#### HELP. kubectl-dynamicns delete --help
`kubectl-dynamicns delete` - Удалить принудительно `dynamic namespace` не дожидаясь **garbage collector**
```
Delete dynamic namespace

Usage:
  kubectl-dynamicns delete [NAMESPACE]  [flags]

Examples:
dynamictl delete dynamic-namespace-dkwcc

Flags:
  -h, --help   help for delete

Global Flags:
   -c, --cluster string    Cluster ENV. Environment CLUSTER (default environment "CLUSTER")
   -P, --password string   Password from LDAP. Environment LDAP_PASSWORD (default environment "LDAP_PASSWORD")
   -U, --user string       User from LDAP. Environment LDAP_USER (default environment "LDAP_USER")

```
#### Example:
##### Short version:
```
export CLUSTER=dev2
export LDAP_USER=v.a.ryumin
export LDAP_PASSWORD=P@ssword
./kubectl-dynamicns delete dynamic-namespace-rfj72

{
  "delete": true
}

```
##### Full version:
```
./kubectl-dynamicns --cluster=dev2 --user=v.a.ryumin --password=P@ssword delete dynamic-namespace-rfj72

{
  "delete": true
}
```

## Errors
|Error message   |Description   | 
|---|---|
|Cannot parse time::: ...  | Это значит что вы передали дату удаление не в формате RFC3339   | 
|Forbidden creating a dynamic namespace. Limit dynamic namespace per user: ...   | Достигнут лимит создаваемых namespace-ов на пользователя **(лимит задается на сервере)**   | 
|Forbidden use more Limits CPU then: ...  | Достигнут лимит по CPU на namespace **(лимит задается на сервере)**   | 
|Forbidden use more Limits Memory then: ...   | Достигнут лимит по Memory на namespace **(лимит задается на сервере)**    | 
