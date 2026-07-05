> 原文: https://turbo.hotwired.dev/handbook/native
> タイトル: Go Native on iOS & Android
> 最終取得: 2026-07-05

# Go Native on iOS & Android

Hotwire Native for iOS は、Turbo に対応した Web アプリをネイティブ iOS シェルで包むためのツール一式を提供します。単一の WKWebView インスタンスを複数の view controller にまたがって管理し、Turbo によるクライアントサイドのパフォーマンス上の恩恵をすべて活かしたまま、ネイティブのナビゲーション UI を実現します。詳細は [Hotwire Native: iOS](https://github.com/hotwired/hotwire-native-ios/) を参照してください。

Hotwire Native for Android も同種のツールを提供し、単一の WebView インスタンスを複数の Fragment 遷移先にまたがって管理します。詳細は [Hotwire Native: Android](https://github.com/hotwired/hotwire-native-android/) を参照してください。

ネイティブアダプタでできることを把握するには、デモ用のネイティブアプリケーションをセットアップするのが最も手っ取り早いです。[iOS 用](https://native.hotwired.dev/ios/getting-started)と [Android 用](https://native.hotwired.dev/android/getting-started)を用意しています。ネイティブ環境でコードを開き、手順に沿って一通りの機能を試してみてください。

[次へ: Building Your Turbo Application](07_building.md)
