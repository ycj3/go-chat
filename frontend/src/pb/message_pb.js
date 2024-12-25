import * as $protobuf from "protobufjs/minimal";

// Common aliases
const $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
const $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

export const chat = $root.chat = (() => {

    /**
     * Namespace chat.
     * @exports chat
     * @namespace
     */
    const chat = {};

    chat.ChatMessage = (function() {

        /**
         * Properties of a ChatMessage.
         * @memberof chat
         * @interface IChatMessage
         * @property {string|null} [user] ChatMessage user
         * @property {string|null} [message] ChatMessage message
         */

        /**
         * Constructs a new ChatMessage.
         * @memberof chat
         * @classdesc Represents a ChatMessage.
         * @implements IChatMessage
         * @constructor
         * @param {chat.IChatMessage=} [properties] Properties to set
         */
        function ChatMessage(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * ChatMessage user.
         * @member {string} user
         * @memberof chat.ChatMessage
         * @instance
         */
        ChatMessage.prototype.user = "";

        /**
         * ChatMessage message.
         * @member {string} message
         * @memberof chat.ChatMessage
         * @instance
         */
        ChatMessage.prototype.message = "";

        /**
         * Creates a new ChatMessage instance using the specified properties.
         * @function create
         * @memberof chat.ChatMessage
         * @static
         * @param {chat.IChatMessage=} [properties] Properties to set
         * @returns {chat.ChatMessage} ChatMessage instance
         */
        ChatMessage.create = function create(properties) {
            return new ChatMessage(properties);
        };

        /**
         * Encodes the specified ChatMessage message. Does not implicitly {@link chat.ChatMessage.verify|verify} messages.
         * @function encode
         * @memberof chat.ChatMessage
         * @static
         * @param {chat.IChatMessage} message ChatMessage message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatMessage.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.user != null && Object.hasOwnProperty.call(message, "user"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.user);
            if (message.message != null && Object.hasOwnProperty.call(message, "message"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.message);
            return writer;
        };

        /**
         * Encodes the specified ChatMessage message, length delimited. Does not implicitly {@link chat.ChatMessage.verify|verify} messages.
         * @function encodeDelimited
         * @memberof chat.ChatMessage
         * @static
         * @param {chat.IChatMessage} message ChatMessage message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatMessage.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a ChatMessage message from the specified reader or buffer.
         * @function decode
         * @memberof chat.ChatMessage
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {chat.ChatMessage} ChatMessage
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatMessage.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.chat.ChatMessage();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.user = reader.string();
                        break;
                    }
                case 2: {
                        message.message = reader.string();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a ChatMessage message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof chat.ChatMessage
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {chat.ChatMessage} ChatMessage
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatMessage.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a ChatMessage message.
         * @function verify
         * @memberof chat.ChatMessage
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        ChatMessage.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.user != null && message.hasOwnProperty("user"))
                if (!$util.isString(message.user))
                    return "user: string expected";
            if (message.message != null && message.hasOwnProperty("message"))
                if (!$util.isString(message.message))
                    return "message: string expected";
            return null;
        };

        /**
         * Creates a ChatMessage message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof chat.ChatMessage
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {chat.ChatMessage} ChatMessage
         */
        ChatMessage.fromObject = function fromObject(object) {
            if (object instanceof $root.chat.ChatMessage)
                return object;
            let message = new $root.chat.ChatMessage();
            if (object.user != null)
                message.user = String(object.user);
            if (object.message != null)
                message.message = String(object.message);
            return message;
        };

        /**
         * Creates a plain object from a ChatMessage message. Also converts values to other types if specified.
         * @function toObject
         * @memberof chat.ChatMessage
         * @static
         * @param {chat.ChatMessage} message ChatMessage
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        ChatMessage.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.user = "";
                object.message = "";
            }
            if (message.user != null && message.hasOwnProperty("user"))
                object.user = message.user;
            if (message.message != null && message.hasOwnProperty("message"))
                object.message = message.message;
            return object;
        };

        /**
         * Converts this ChatMessage to JSON.
         * @function toJSON
         * @memberof chat.ChatMessage
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        ChatMessage.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        /**
         * Gets the default type url for ChatMessage
         * @function getTypeUrl
         * @memberof chat.ChatMessage
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        ChatMessage.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/chat.ChatMessage";
        };

        return ChatMessage;
    })();

    return chat;
})();

export { $root as default };