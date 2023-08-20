# Artifact Store

## The idea

You know the file structure of your bundles for your app, thus you can have defined types for your bundles as well as the metadata of the bundle that you have defined.

The index.yaml contains user defined metadata for you to be able to understand what is in your store. So you can query it, in one go by reading a unique file


## Config controlled by application

## Store structure

### File structure at the store root

```bash
artifact-store/
├─ index.yaml
├─ bucket1 (eg: pizzahut)/             // a bucket is a type (a struct or aliasies to give name)
│  ├─ bundle1 (eg: default & v1.0.0)   // a bundle has a unique composed key of name and version
│  ├─ bundle2 (eg: default & v1.0.1)
│  ├─ bundle3 (eg: mega    & v1.0.0)
│  ├─ ...
├─ ...
```

### File structure of a bundle

```bash
bundle_x/
├─ ./bundle             // Metafiles for the store system
│  ├─ MANIFEST.yaml
│  ├─ sha256sum.txt
│  ├─ ...
├─ file_1.txt
├─ file_2.txt
├─ file_3.txt
├─ ...
```

### Inside the MANIFEST.yaml

```yaml
name: ""         // name of the bundle
version: ""      // version of the bundle
metadata:
  ...            // abritrary yaml
```

### index.yaml

- The index have the list of all the MANIFEST.yaml in every bundles stored as a list.
- The list is divided by buckets.
- Each entry can have more information then just the MANIFEST.yaml to complete with what the store needs to manage the bundles.
  - For example, like the url of the tarball, the sha1sum, url, etc...

```yaml
apiVersion: v1alpha
entries:
  manifests:
    - version: v1.0.0
      name: tmtc
      metadata:
        name: "test"
        description: "second test"
        tags:
          vehicule: neutron
          application: muon
      url: ./tmtc_v1.0.0.tgz
      files:
        - .bundle/
        - .bundle/MANIFEST.yaml
        - .bundle/sha1sums.txt
        - tmtc.bfbs
  pizzahut:
    - version: v1.0.0
      name: default
      metadata:
        name: "general_store"
        description: "configuration for defacto store"
        tags:
          general: true
          store: pizza
      url: ./default_v1.0.0.tgz
      files:
        - .bundle/
        - .bundle/MANIFEST.yaml
        - .bundle/sha1sums.txt
        - pizzahut.yaml
```

### To bundle

```bash
store bundle <name> <version> <metadata> -- <files>
```

- a possible idea, to have a own app bundle command, the metadata has defaults, or can show an example, or validate.
- could be even extended to the files inside the bundle
- could have a generate command, that generate you a fresh bundle with all the defaults.
- using generic cmd

### To upload

```bash
store upload <bucket> <path_to_bundled_bundle> -t <store_server_url>
```

### Versionning of metadata in bundle

the dev maintains the metadata in the bundle, by defining its own struct and managing it himself like so

```go
type ConfigHeader[T any] struct {
  ApiVersion  string "json:'apiVersion'"
  Metadata T      "json:'metadata'"
}
```

this could be offered as a ready made struct

### Tags

use a master struct that has all the configs type, use tags to control to which files


## Idea to explore

### Builder pattern

- with_store
  - with_local_disk
  - with_s3



### Utility 'load all' functions

as your own program start up, you have a set of utility functions to load your bundles in memory.

- have a search query based on name, version, your own metadata, or simply all.
- that returns a list of bundles that match those criteria.
- then you could do: `get_file_content[<your_type>](list_of_bundles, <file_path_in_bundle>) []T`
- instead of generic, we can also use `get_file_content(list_of_bundles, <file_path_in_bundle>, &var []YOUR_TYPE)
- or have maybe a struct like this

```go
type BundleData struct {   // Arbitrary struct defined by the dev
  File1 PizzaHut           "storePath='<file_path_in_bundle>'" 
  File2 ThreadConfig       "storePath='<file_path_in_bundle>'"
  File3 LogConfig          "storePath='<file_path_in_bundle>'"
}
```

- we could defined our own type and have custom tags that point where the file is in the bundle.
- you could have functions that load all the bundles of a type
- if it was a store connected to a parent proxy, the load all functions could be applied to its local store only or both local and remote.
- with remote, it could dangerous with size if the search query would be too broad.
- the best is it stays local when it was asked to have a copy of it.

### Proxy of a store

- a store could be running as a stand alone thing and manage its own store with local or s3 behind.
- lower store or store running along side app, could connect to a parent store and proxy the request to it.
- idea: use some voodo kubernetes magic to connect everything together
- idea: uploading a bundle to a local store that his connected to a remote store, will upload to the remote store as well.
- idea: it could be also a fully distributed system. some store could have different settings such as longer ttl of bundles, etc

### More commands to manage the store

You could have Create/Delete/Move/Copy/Rename and other utility function to manage the store once the files are there, create buckets, etc.

### Bucket created by application

- In code, you could do, `create bucket for type` and it would create a bucket for the type, and the bucket would be the name of the type.
- so if other app needs the same type of bundle, they dont need to know the name of the bucket, they just need to know the type.
- you could even maintain version with type aliases.

## Alternate option, There is no metadata in the bundle

What changes

### Inside the MANIFEST.yaml

- No arbritraty metadata

```yaml
name: ""         // name of the bundle
version: ""      // version of the bundle
```

### index.yaml

- No arbritraty metadata

```yaml
apiVersion: v1alpha
entries:
  manifests:
    - version: v1.0.0
      name: tmtc
      url: ./tmtc_v1.0.0.tgz
      files:
        - .bundle/
        - .bundle/MANIFEST.yaml
        - .bundle/sha1sums.txt
        - tmtc.bfbs
  pizzahut:
    - version: v1.0.0
      name: default
      url: ./default_v1.0.0.tgz
      files:
        - .bundle/
        - .bundle/MANIFEST.yaml
        - .bundle/sha1sums.txt
        - pizzahut.yaml
```

### To bundle

to bunle, we remove the metadata file arg

```bash
store bundle <name> <version> -- <files>
```

### Versionning of metadata in bundle

the dev maintains the metadata in the bundle by himself using a file that in reserve for it in the bundle. 

```bash
bundle_x/
├─ ./bundle             // Metafiles for the store system
│  ├─ MANIFEST.yaml
│  ├─ sha256sum.txt
│  ├─ ...
├─ METADATA.yaml        // It's own metadata file maintained by the dev
├─ file_2.txt
├─ file_3.txt
├─ ...
```

- Technically, this options is less tied up by lot of typing and more freedom, but it's then more difficult to have search queries and understand what its your store so you can manipulate it.
- Moreover, the speed of reading one file is faster then reading a lake of bundle.
- And technically, using the original idea, you can just have no metadata.
- And you dont need to use the typing, you could just define everything as []byte, and then manipulate that directly to your own struct.
