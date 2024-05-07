from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ScanRequest(_message.Message):
    __slots__ = ("content",)
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    content: str
    def __init__(self, content: _Optional[str] = ...) -> None: ...

class ScanResponse(_message.Message):
    __slots__ = ("result",)
    class Result(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        UNKNOWN: _ClassVar[ScanResponse.Result]
        HAM: _ClassVar[ScanResponse.Result]
        SPAM: _ClassVar[ScanResponse.Result]
    UNKNOWN: ScanResponse.Result
    HAM: ScanResponse.Result
    SPAM: ScanResponse.Result
    RESULT_FIELD_NUMBER: _ClassVar[int]
    result: ScanResponse.Result
    def __init__(self, result: _Optional[_Union[ScanResponse.Result, str]] = ...) -> None: ...
