![](../../img/providers/snyk.png)

# Snyk

In order to use `bomber` with Snyk you will need to be a Snyk customer. Access requires your Snyk API Token, which you can retrieve from the web interface or by running: 

```
snyk config get api
```

Once you have your token you can run bomber like so: 

```
bomber scan --provider snyk --token xxx sbom.json
```

Note rather than passing the API token explicitly, you can also set this as an environment variable, either as `SNYK_TOKEN` or the generic `BOMBER_PROVIDER_TOKEN`.

By default, `bomber` will use Snyk's global API (https://api.snyk.io). To use a different Snyk API, you can specify its base URL on the `SNYK_API` environment variable.

```
SNYK_API=https://api.eu.snyk.io bomber scan --provider snyk sbom.json
```


## Supported ecosystems

At this time, the Snyk provider supports the following ecosystems:

* npm
* Maven
* CocoaPods
* Composer
* RubyGems
* Nuget
* PyPi
* Hex
* Cargo
* Swift
* C/C++
* apk
* Debian
* Docker
* RPM
