import re
from typing import Dict, Optional, Union

RequestDict = Dict[str, Union[str, Dict[str,str]]]

def parse_raw_request(
    request: bytes,
    encoding: str = 'utf-8'
) -> RequestDict:
    pattern: re.Pattern = re.compile(
        pattern=r'(?P<method>\w+)\s+(?P<path>[^\s?]+)\s+HTTP/(?P<version>\d+\.\d+)\r\n(?P<headers>(?:.*\r\n)*?)\r\n(?P<body>.*)',
        flags=re.DOTALL
    )
    
    match: Optional[re.Match] = pattern.match(
        request.decode(encoding=encoding)
    )

    if not match:
        raise ValueError('Invalid HTTP request')
    return {
        'method': match.group('method'),
        'path': match.group('path'),
        'version': match.group('version'),
        'headers': dict(re.findall(r'(?P<name>.*?): (?P<value>.*?)\r\n', match.group('headers'))),
        'body': match.group('body')
    }