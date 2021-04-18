"""Utility for ICU."""

import threading
from typing import List, Union


try:
    import icu
except ImportError:
    msg = (
        f"`icu` needs to be installed for {__file__}.\n"
        "It's recommended to install `icu` using `conda`.\n"
        "    conda install -y icu\n"
        "    CFLAGS='-std=c++11' pip install PyICU"
    )
    raise ImportError(msg)


_mutex = threading.Lock()
_break_iterator = icu.BreakIterator.createWordInstance(icu.Locale.getUS())
# US locale seem to work well for almost all languages.


def icu_break(text: Union[icu.UnicodeString, str]) -> List[int]:
    """Get the breakpoints of a given text using icu word break iterator.

    Args:
        text (Union[icu.UnicodeString, str]): Input text.

    Returns:
        List[int]: Breakpoints except for 0.
    """
    with _mutex:
        _break_iterator.setText(text)
        return list(_break_iterator)


def icu_tokenize(text: Union[icu.UnicodeString, str]) -> List[str]:
    """Tokenize text using icu word break iterator.

    Args:
        text (Union[icu.UnicodeString, str]): Input text.

    Returns:
        List[str]: Tokens.
    """
    if isinstance(text, str):
        text = icu.UnicodeString(text)

    segments = []
    p0 = 0
    for p1 in icu_break(text):
        segment = text[p0:p1].strip()
        if len(segment) > 0:
            segments.append(str(segment))
        p0 = p1
    return segments
