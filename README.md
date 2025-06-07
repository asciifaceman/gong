# Gong
A simple asset packer in pure go inspired by thehugh100/Simple-Asset-Packer

# Headers

### Container Header

`gonginfo_t`
| Field | Length | Code | Description |
| ----- | ------ | ---- | ----------- |
| SOG   | 2      | FF D6| Start of Gong |
| IDENT | 5      | 47 4F 4E 47 00 | "GONG" null terminated |
| ...   |        |      | content |
| EOG   | 2      | FF D7 | End of gong |

### Container Table





### Asset Header

`gongasset_t`
| Field | Length | Code | Description |
| ----- | ------ | ---- | ----------- |
| SOA   | 2      | FF D8| Start of Asset |
| ASST  | 4      | 41 53 53 54 | Asset leading header |
| FNL | 4 | | File name length |
| FN | n | | File name |
| LEN   | 4      | 00 00 00 00 | int content length |
| CON   | n      | content     | content of the asset (NULL terminated) |
| EOA   | 2      | FF D9       | End of asset |


# CLI Usage

* list
  * list assets in bundle
* add
  * add asset to bundle
* remove
  * remove asset from bundle (by name)
* 