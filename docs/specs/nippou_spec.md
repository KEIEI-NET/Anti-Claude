# プログラム仕様書: Nippou Domain Entity

## 1. 概要

日報の核心的なビジネスルールとデータ構造を定義するエンティティ。純粋なドメイン知識のみを持ち、外部ライブラリへの依存を最小限（`uuid`等のみ）とする。

## 2. データ構造 (Value Objects)

実装段階で、不変条件を保証するために以下の型安全なバリューオブジェクトを導入。

- **ID**: `string` ベースのユニーク識別子。
- **Location**: `Latitude`, `Longitude`, `Address` を持つ不変構造体。
- **VoiceConfig**: 音声の有効化状態と使用モデル名のペア。
- **Tag**: 1文字以上のバリデーション済み文字列。

## 3. 生成ロジック (Builder Pattern)

複雑な初期化を安全に行うため、`NippouBuilder` を採用。

- `NewNippouBuilder(date, content)`: 必須項目で初期化。
- `WithLocation(loc)` / `WithVoice(v)` / `WithTags(tags)`: オプション項目。
- `Build()`: 最終的な不変条件チェックを実行。

## 4. ビジネスルール (Invariants)

- **Content**: 空文字不可、最大10,000文字。
- **Date**: YYYY-MM-DD フォーマット厳守。
- **Location**: 緯度(-90 to 90), 経度(-180 to 180)。
- **Tags**: 重複不可、一要素当たりの最大長、最大数制限。

## 5. リポジトリ・インターフェース

`Repository` インターフェースを定義し、Infrastructure層での具象実装を要求する。

- `Save(n *Nippou) error`
- `FindByID(id ID) (*Nippou, error)`
- `FindByDate(date time.Time) ([]*Nippou, error)`
