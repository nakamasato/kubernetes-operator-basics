# Run Operator

## Overview

We'll run **PostgreSQL** with and without [postgres-operator](https://github.com/zalando/postgres-operator)

## 1. Run PostgreSQL without operator.

1. Create Yaml file.
1. Apply

    ```
    kubectl delete -f postgres-sts.yaml
    ```

## 2. Run PostgreSQL with operator.


1. Install Postgres Operator

    ```
    kubectl apply -k github.com/zalando/postgres-operator/manifests
    ```
1. Deploy the Operator UI

    ```
    kubectl apply -k github.com/zalando/postgres-operator/ui/manifests
    ```

    Check:

    ```
    kubectl port-forward svc/postgres-operator-ui 8081:80
    ```

    Open http://localhost:8081/

1. Create a Postgres Cluster


    ```
    kubectl create -f https://raw.githubusercontent.com/zalando/postgres-operator/master/manifests/minimal-postgres-manifest.yaml
    ```

    <details><summary>Roles and Databases initially created:</summary>

    yaml:

    ```yaml
      users:
        zalando:  # database owner
        - superuser
        - createdb
        foo_user: []  # role for application foo
      databases:
        foo: zalando  # dbname: owner
      preparedDatabases:
        bar: {}
    ```

    roles:

    ```
    \du
                                                         List of roles
        Role name    |                         Attributes                         |               Member of
    -----------------+------------------------------------------------------------+----------------------------------------
     admin           | Create DB, Cannot login                                    | {foo_user,zalando,bar_owner}
     bar_data_owner  | Cannot login                                               | {bar_data_writer,bar_data_reader}
     bar_data_reader | Cannot login                                               | {}
     bar_data_writer | Cannot login                                               | {bar_data_reader}
     bar_owner       | Cannot login                                               | {bar_writer,bar_data_owner,bar_reader}
     bar_reader      | Cannot login                                               | {}
     bar_writer      | Cannot login                                               | {bar_reader}
     foo_user        |                                                            | {}
     postgres        | Superuser, Create role, Create DB, Replication, Bypass RLS | {}
     robot_zmon      | Cannot login                                               | {}
     standby         | Replication                                                | {}
     zalando         | Superuser, Create DB                                       | {}
     zalandos        | Cannot login                                               | {}
    ```

    databases:

    ```
    \l
                                      List of databases
       Name    |   Owner   | Encoding |   Collate   |    Ctype    |   Access privileges
    -----------+-----------+----------+-------------+-------------+-----------------------
     bar       | bar_owner | UTF8     | en_US.utf-8 | en_US.utf-8 |
     foo       | zalando   | UTF8     | en_US.utf-8 | en_US.utf-8 |
     postgres  | postgres  | UTF8     | en_US.utf-8 | en_US.utf-8 |
     template0 | postgres  | UTF8     | en_US.utf-8 | en_US.utf-8 | =c/postgres          +
               |           |          |             |             | postgres=CTc/postgres
     template1 | postgres  | UTF8     | en_US.utf-8 | en_US.utf-8 | =c/postgres          +
               |           |          |             |             | postgres=CTc/postgres
    (5 rows)
    ```

    </details>

1. Check

    ```
    # check the deployed cluster
    kubectl get postgresql

    # check created database pods
    kubectl get pods -l application=spilo -L spilo-role

    # check created service resources
    kubectl get svc -l application=spilo -L spilo-role
    ```

1. Connect to Postgres cluster.

    ```
    kubectl exec -it acid-minimal-cluster-0 -- psql -Upostgres
    psql (14.0 (Ubuntu 14.0-1.pgdg18.04+1))
    Type "help" for help.

    postgres=#
    ```

1. Check on UI.

    http://localhost:8081/#/status/default/acid-minimal-cluster

1. Clean up.

    ```
    kubectl delete -f https://raw.githubusercontent.com/zalando/postgres-operator/master/manifests/minimal-postgres-manifest.yaml
    kubectl delete -k github.com/zalando/postgres-operator/ui/manifests
    kubectl delete -k github.com/zalando/postgres-operator/manifests
    ```

## Other Postgres Operator
- https://www.kubegres.io/
