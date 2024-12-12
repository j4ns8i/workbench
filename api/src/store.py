import xxhash


_topic = "topic"
_messages = "messages"


def topic_key(topic: str) -> str:
    """
    Generate a formatted store key for a given topic.

    Args:
        topic (str): The topic to hash.

    Returns:
        str: The formatted store key.
    """
    return f"{_topic}:{xxhash.xxh64_hexdigest(topic)}:{_messages}"
