# プログラム仕様書: Nippou Create UseCase (Application Layer)

## 1. 概要

`Nippou` ドメインモデルをオーケストレートし、入出力バリデーション、エンティティの構築、永続化（Repository）の呼び出しを統括する。

## 2. モジュール構成 (File Splits)

実装段階で、保守性と再利用性を高めるために以下の3ファイルに分割。

- **errors.go**: アプリケーション層固有の `UseCaseError` 体系。
- **dto.go**: `CreateInput` (Request) と `CreateOutput` (Response) の定義。
- **interactor.go**: ビジネスロジックの本体実装。

## 3. 詳細仕様: interactor.go

### 3.1 実行フロー (Execute)

1. **Context Check**: コンテキストのキャンセルを早期に判定。
2. **DTO Validation**: `CreateInput.Validate()` による早期バリデーション。
3. **Entity Build**:
    - `Domain.NewNippouBuilder` を利用。
    - Location, VoiceConfig, Tags の順次付与（ドメインルール違反時は `DOMAIN_VIOLATION` エラー）。
4. **Persistence**: `Repository.Save` の実行。
5. **Output Mapping**: 永続化成功後、ドメインモデルを `CreateOutput` DTO へ変換。

## 4. エラーハンドリング (errors.go)

以下のコード体系により、呼び出し元（Interface層）が適切な HTTP ステータスを判断できるようにする。

- `INVALID_INPUT`: クライアント側の入力不備。
- `REPOSITORY_ERROR`: 永続化失敗。
- `DOMAIN_VIOLATION`: ドメイン知識に反する操作。
- `CONTEXT_CANCELLED`: タイムアウト。

## 5. テスト仕様 (interactor_test.go)

44件のテストケースにより、以下の境界条件を網羅。

- **Success Cases**: 全オプションあり、一部欠落、境界値（緯度90.0等）。
- **Error Cases**: 重複タグ、不正な日付フォーマット、コンテキストキャンセル。
- **Concurrency**: 並列実行時の安全性。
