# Artifact Management

This artifact platform is used to store any artifacts needed by our applications. This platform use a strict system on how to package files to be consistent accross all artifacts. Please refer to the Artifact file structure at the bottom.

## MaxGds TmTc
To push max-tmtc (cmd_tlm_db) manually to the artifact platform, you must package the cmd_tlm_db and then push it to the s3. Two scripts are provided:
 - `package.sh` > Package files in the proper format.
 - `s3cmd.sh` > Push or retrieve artifact.

```sh
./package.sh <package> <version> -- <file1> <file2> <file3> <file4> <file5> <file6> ...

./s3cmd.sh <operation> s3://<bucket>/<release>/<package>/<artifact> <path_to_artifact>
```

Here's an example on how to deploy a maxgds max-tmtc
```sh
./package.sh max-tmtc empty-1.0.0 -- FswCmdDetail.csv FswCmdEnum.csv FswCmdMaster.csv FswTlmDetail.csv FswTlmEnum.csv FswTlmMaster.csv

./s3cmd.sh put s3://pr-sw-artifactrepo-01/repeater/max_tmtc/max-tmtc_empty-1.0.0.tgz max-tmtc_empty-1.0.0.tgz
```

---
## Artifact File Structure

URL:
- `s3://[bucket_environment]/[release|branch]/[component]/artifact-[buildid|version].tgz`
- `s3://dv-sw-artifactrepo-01/LunarD/max-tmtc/max-tmtc-1.0.0.tgz`
- `s3://st-sw-artifactrepo-01/develop/max-tmtc/max-tmtc-1.0.0.tgz`
- `s3://pr-sw-artifactrepo-01/local/max-tmtc/max-tmtc-1.0.0.tgz`


## ARTIFACT:
```
--- index.yaml (later)
  + develop
    + max-tmtc
    | + max-tmtc-1.0.0.tgz
    | |+ .bundle
	| |  |- sha1sums.txt
	| |  |- MANIFEST.yaml
    | |- XXX.csv
    | |- XXX.csv
    | |- XXX.csv
    | + max-tmtc-1.0.0.tgz.prov (later)
    + covalence
    |-+ covalence-1.0.0.tgz
    | |+ .bundle
	| |  |- sha1sums.txt
	| |  |- MANIFEST.yaml
      |- tmtc
```

Steps to get artifact (client):
1. Get release+package+version
2. curl url
3. tar -x in working directory

Steps for lookup (client):
1. Bucket comes from env var.
2. One dropdown for [release|branch]
3. One dropdown for the rest using prefix of s3 `curl  -GET flashblade.com?prefix=artifactRepoDev/R3.5/max-tmtc`


Index.yaml keep an index of all artifacts for easly look up and maintenance of our artifact system.

index.yaml (to be refined)
--------------
repo: develop
releases:
- name: develop
  description: Develop builds
- name: R3.5
  description: Photon 3.5 Release

Artifacts can be multiple objects. They will always be archive into one single objects with an extra .bundle folder containing MANIFEST.yaml, sha1sums, and other file that provides extra informations.
In the end, each artifacts is only one file.

MANIFEST.yaml:
--------------
```yaml
name: tlm_db
version: 1.0.0
metadata:
	key1: value
	key2: value
	...
build:
	id:
	sha:
	branch:
	repo:
```

Inspiration
-----------

https://v2.helm.sh/docs/chart_repository/#the-chart-repository-structure

INDEX.YAML EXAMPLE IN https://muck.rl.team/repository/sigmanet/index.yaml

```yaml
apiVersion: v1
entries:
  sigmanet-site:
  - appVersion: 0.0.0
    created: 2021-05-06T03:50:34.994Z
    description: A Helm chart for Kubernetes
    digest: 0323c17cce1b80cb7f6c901f4d99abd0901526b157f0fb1ec53d6bf306aa4e14
    name: sigmanet-site
    urls:
    - sigmanet-site-0.1.0.tgz
    version: 0.1.0
  rocketdb:
  - appVersion: 2.0.528
    created: 2021-03-18T21:04:04.444Z
    description: SigmaNet TSDB
    digest: 07c117bf877f5b2936dd3a27a2c84b3633226ae596148ec74be8dd2823c4778b
    name: rocketdb
    urls:
    - rocketdb-0.1.0-dev.tgz
    version: 0.1.0-dev
  - appVersion: 2.0.563
    created: 2021-03-30T01:13:22.384Z
    description: SigmaNet TSDB
    digest: 8916e5992385d75e897532ac49ec1fff6f1212f9b1bc4ee522b9a59388924583
    name: rocketdb
    urls:
    - rocketdb-0.1.1-dev.tgz
    version: 0.1.1-dev
  - appVersion: 1.0.310
    created: 2021-05-12T21:43:44.880Z
    description: SigmaNet TSDB
    digest: 5d8343bfbdf0798b346f7343e514def4938b0781da217267e15321d5d8c96532
    name: rocketdb
    urls:
    - rocketdb-1.0.310.tgz
    version: 1.0.310
  - appVersion: 2.0.726
    created: 2021-05-06T03:39:12.403Z
    description: SigmaNet TSDB
    digest: 7aef67624244e595f079402eebe165c4ffa5373d54af689ed4fbafb46a28323e
    name: rocketdb
    urls:
    - rocketdb-2.0.726.tgz
    version: 2.0.726
generated: 2021-08-31T04:04:09.985Z
```
