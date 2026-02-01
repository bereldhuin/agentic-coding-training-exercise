import Foundation

/// Cursor codec for encoding/decoding pagination cursors
enum CursorCodec {
    /// Encode cursor data to base64 JSON string
    static func encode(id: Int, createdAt: String) -> String {
        let cursorData = CursorData(id: id, createdAt: createdAt)
        guard let jsonData = try? JSONEncoder().encode(cursorData),
              let jsonString = String(data: jsonData, encoding: .utf8) else {
            return ""
        }
        return jsonString.data(using: .utf8)?.base64EncodedString() ?? ""
    }

    /// Decode base64 JSON string to cursor data
    static func decode(_ cursor: String) -> CursorData? {
        guard let data = Data(base64Encoded: cursor),
              let cursorData = try? JSONDecoder().decode(CursorData.self, from: data) else {
            return nil
        }
        return cursorData
    }
}
