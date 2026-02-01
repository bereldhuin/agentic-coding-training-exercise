import Foundation

/// ISO8601 date formatter for consistent date encoding/decoding
enum ISO8601Formatter {
    static let shared: ISO8601DateFormatter = {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [
            .withInternetDateTime,
            .withFractionalSeconds
        ]
        return formatter
    }()

    /// Encode Date to ISO8601 string
    static func string(from date: Date) -> String {
        return shared.string(from: date)
    }

    /// Decode ISO8601 string to Date
    static func date(from string: String) -> Date? {
        return shared.date(from: string)
    }
}

/// Custom date encoding strategy for JSONEncoder/JSONDecoder
extension JSONEncoder {
    static var vaporISO8601: JSONEncoder {
        let encoder = JSONEncoder()
        encoder.dateEncodingStrategy = .custom { date, encoder in
            var container = encoder.singleValueContainer()
            try container.encode(ISO8601Formatter.string(from: date))
        }
        encoder.outputFormatting = [.sortedKeys, .withoutEscapingSlashes]
        return encoder
    }
}

extension JSONDecoder {
    static var vaporISO8601: JSONDecoder {
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .custom { decoder in
            let container = try decoder.singleValueContainer()
            let string = try container.decode(String.self)
            guard let date = ISO8601Formatter.date(from: string) else {
                throw DecodingError.dataCorruptedError(
                    in: container,
                    debugDescription: "Expected date string to be ISO8601-formatted"
                )
            }
            return date
        }
        return decoder
    }
}
