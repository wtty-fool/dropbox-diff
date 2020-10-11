# dropbox-diff
> Check what is missing in your Dropbox folder the easy way.

`dropbox-diff` produces a list of differences between specified Dropbox directory
and local directory, without a need to download the files first.

## Usage
```
drobox-diff --dropbox /pictures/mountains ~/pictures/trip-to-vermont
```

Requires OAuth2 access token in `token` file. You can generate one at https://www.dropbox.com/developers/apps/.
