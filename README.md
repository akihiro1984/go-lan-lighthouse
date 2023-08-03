# GO LAN Lighthouse

## Overview

IoT 端末が DHCP が有効なネットワーク上を経由している場合に
再起動時点で IP が変わってしまい、行方不明になるので捜索の為のソリューション。

理由があって固定 IP を振ることが不可能なネットワークにもマイクロサービスを蔓延らせる事が可能である。

## System

### Tower Mode

通信可能 IP をレスポンスするサーバモード

Agent の接続先、情報 HUB のターゲットサーバ等で稼働させておくことを想定。

### Ship Mode

IoT端末等、訳あって電源が入切が著しい端末で動くモード。

稼働時にブロードキャストパケットを送信し、Towerからの IP レスポンスを待つ。
IPが判明したときには終了する。

## Sample

Sample ディレクトリに動作サンプルを格納した。
シェルスクリプトを動作させると、.envファイルにIPを格納する。

## Contributing

必要最小限の機能しか実装していない為、環境ごとに不具合があるかもしれない。
また機能追加等のリクエスト大歓迎。

また機能追加のリクエストは、以下のプルリクエスト手法に基づくものとする。
https://github.com/MarcDiethelm/contributing/blob/master/README.md


## Todo

- ネットワークカードが複数ある場合の IP 選定方法

## Disclaimer

拙作プログラムの利用に関しては、自己責任でお願い申し上げます。

## Acknowledgments

[Howto: UDP broadcast in Go](https://github.com/aler9/howto-udp-broadcast-golang)

[Help With UDP Broadcast](https://forum.golangbridge.org/t/help-with-udp-broadcast/22036)


