# Sources for the names and their explanations

## Editorial scope

The existing 99 Arabic identities, numbering and order are retained. The short
English labels are concise glosses; descriptions and summaries expand the same
meaning. Closely related names can have overlapping meanings. The compound
Dhu Al Jalaali Wa Al Ikraam uses Majestic and Generous as its short label; the
description and summary explain possession of majesty, honour and generosity.

This is the familiar enumeration, not a claim that all 99 listed forms occur
literally in the Quran or that scholars agree on one exhaustive enumeration.
The enumerated narration in [Tirmidhi 3507](https://sunnah.com/tirmidhi:3507)
is graded weak. A related Quranic action alone does not authenticate every noun
that can be derived from it as a divine name.

## Sources and interpretive choices

The primary source for the Quran references is the Arabic verse itself in
`quran/data/arabic.json`. The table below links every retained reference. The
reference note distinguishes nominal forms, inflection, plurals, verbs and
related passages. This is a curated reference list, not a concordance of every
occurrence. New English explanations are editorial paraphrases, not translations
quoted from a particular published edition.

Additional interpretation and hadith sources:

- [Ibn Kathir on 59:23–24](https://quran.com/59:23/tafsirs/en-tafisr-ibn-kathir): holiness, peace, security/confirmation of truth, oversight, and the names relating to creation.
- [Ibn Kathir on 4:85](https://quran.com/4:85/tafsirs/en-tafisr-ibn-kathir): Muqeet includes interpretations of nourishment, power and watchfulness; do not force its meaning into a rigid distinction from Razzaaq.
- [Ibn Kathir on 112:1–4](https://quran.com/112:2/tafsirs/en-tafisr-ibn-kathir) and [Ma'arif al-Qur'an on 112:2](https://quran.com/112:2/tafsirs/en-tafsir-maarif-ul-quran): Ahad and Samad. The agreed short label for Samad is The Eternal Refuge, with needlessness and universal dependence explained in the description and summary.
- [Al-Ghazali, Al-Maqsad al-Asna, original Arabic](https://www.ghazali.org/books/asna-bejou.pdf), PDF page 57 (printed 108–109): Muhsi, originating/restoring, giving life/causing death and Qayyoom; PDF page 58 (printed 110–111): Qayyoom, Waajid, Maajid, Waahid, Samad and Qaadir/Muqtadir. Waajid follows the explanation of lacking nothing necessary to perfection, rather than the previous perception summary. This interpretive choice is attributed in the summary. The related Quran reference is 35:15 (Ghaniyy); 38:44 about finding Ayyub patient was removed.
- The earlier descriptions of Mutakabbir, Razzaaq, Kabeer, Majeed/Maajid, Tawwaab and Afuww were checked against the corresponding sections of al-Ghazali's work; [Stade's English translation](https://www.ghazali.org/books/Ghazali-99-Name-of-God.pdf) was used in that earlier pass. Distinctions are explanatory emphases, not claims of mutually exclusive attributes.
- [Abu Dawud 3451](https://sunnah.com/abudawud:3451): Qaabid and Baasit explicitly in the Arabic; Quran 2:245 uses verbs.
- [Muslim 771a](https://sunnah.com/muslim:771a): Muqaddim and Mu’akhkhir explicitly in the Arabic. The English rendering on that page should not be used to identify their literal Arabic forms.
- [Muslim 2713a](https://sunnah.com/muslim:2713a): explanation of Awwal, Aakhir, Zaahir and Baatin. The summaries avoid adding claims that Allah occupies a created physical space.
- [Bukhari 6099](https://sunnah.com/bukhari:6099): divine patience expressed with a comparative form, rather than the literal noun Al Saboor.

## Page composition and updates

`names/data/names.json` supplies Arabic, transliteration, short meaning,
description, Quran references and `reference_note`. The numbered files supply
summaries. Their duplicate name/description fields now match the main records,
with a regression check to prevent drift. The modern detail page, share text and
bookmark excerpts use the corrected summary; the legacy HTML also displays the
description. Quran-reference notes appear in the API, modern page, legacy HTML
and Markdown output. Existing typography and controls are preserved.

Names in the persistent keyword index are replaced transactionally when a hash
of their indexed content changes. The other sources are untouched. Existing
vector matches are resolved back to current Arabic identities before display or
LLM summarisation; their prebuilt semantic ranking is not regenerated. New
keyword searches include the current summaries. Browser search-history entries
containing old names evidence are invalidated once; unrelated history is kept.
The existing content-addressed offline build updates the names corpus.
Previously generated historical daily records and user-saved bookmarks are not
rewritten; new generation and newly saved/shared text use the corrected corpus.

## Current names and Quran references

A blank qualification means the references contain the name or a corresponding
nominal form referring to Allah; it does not require the definite article or the
same grammatical case as the standalone listing.

| # | Name | Short label | Quran references | Qualification |
| --- | --- | --- | --- | --- |
| 1 | Ar Rahmaan | The Beneficent | [1:3](https://quran.com/1/3), [17:110](https://quran.com/17/110), [55:1](https://quran.com/55/1) |  |
| 2 | Ar Raheem | The Merciful | [2:163](https://quran.com/2/163), [3:31](https://quran.com/3/31), [4:100](https://quran.com/4/100), [5:3](https://quran.com/5/3) |  |
| 3 | Al Malik | The King | [20:114](https://quran.com/20/114), [23:116](https://quran.com/23/116), [59:23](https://quran.com/59/23), [62:1](https://quran.com/62/1) |  |
| 4 | Al Quddus | The Holy | [59:23](https://quran.com/59/23), [62:1](https://quran.com/62/1) |  |
| 5 | As Salaam | Source of Peace | [59:23](https://quran.com/59/23) |  |
| 6 | Al Mu'min | Giver of Security | [59:23](https://quran.com/59/23) |  |
| 7 | Al Muhaymin | The Guardian | [59:23](https://quran.com/59/23) |  |
| 8 | Al Azeez | The Almighty | [3:6](https://quran.com/3/6), [4:158](https://quran.com/4/158), [9:40](https://quran.com/9/40), [48:7](https://quran.com/48/7) |  |
| 9 | Al Jabbaar | The Compeller | [59:23](https://quran.com/59/23) |  |
| 10 | Al Mutakabbir | The Supremely Great | [59:23](https://quran.com/59/23) |  |
| 11 | Al Khaaliq | The Creator | [6:102](https://quran.com/6/102), [13:16](https://quran.com/13/16), [39:62](https://quran.com/39/62), [40:62](https://quran.com/40/62), [59:24](https://quran.com/59/24) |  |
| 12 | Al Baari | The Maker | [59:24](https://quran.com/59/24) |  |
| 13 | Al Musawwir | The Fashioner | [59:24](https://quran.com/59/24) |  |
| 14 | Al Ghaffaar | The Forgiving | [20:82](https://quran.com/20/82), [38:66](https://quran.com/38/66), [39:5](https://quran.com/39/5), [40:42](https://quran.com/40/42), [71:10](https://quran.com/71/10) |  |
| 15 | Al Qahhaar | The Subduer | [13:16](https://quran.com/13/16), [14:48](https://quran.com/14/48), [38:65](https://quran.com/38/65), [39:4](https://quran.com/39/4), [40:16](https://quran.com/40/16) |  |
| 16 | Al Wahhaab | The Bestower | [3:8](https://quran.com/3/8), [38:9](https://quran.com/38/9), [38:35](https://quran.com/38/35) |  |
| 17 | Ar Razzaaq | The Provider | [51:58](https://quran.com/51/58) |  |
| 18 | Al Fattaah | The Opener | [34:26](https://quran.com/34/26) |  |
| 19 | Al 'Aleem | The All-Knowing | [2:158](https://quran.com/2/158), [3:92](https://quran.com/3/92), [4:35](https://quran.com/4/35), [24:41](https://quran.com/24/41), [33:40](https://quran.com/33/40) |  |
| 20 | Al Qaabid | The Constrictor | [2:245](https://quran.com/2/245) | Quran 2:245 describes restricting and expanding through verbs. Al Qaabid and Al Baasit are named in Abu Dawud 3451. |
| 21 | Al Baasit | The Expander | [2:245](https://quran.com/2/245) | Quran 2:245 describes restricting and expanding through verbs. Al Qaabid and Al Baasit are named in Abu Dawud 3451. |
| 22 | Al Khaafid | The Abaser | [95:5](https://quran.com/95/5) | Related passage: 95:5 describes bringing low; it does not contain Al Khaafid as a divine name. |
| 23 | Ar Raafi' | The Exalter | [58:11](https://quran.com/58/11), [6:83](https://quran.com/6/83) | Related passages: these verses describe raising in rank through verbs, rather than naming Allah Ar Raafi'. |
| 24 | Al Mu'izz | The Honourer | [3:26](https://quran.com/3/26) | Related passage: 3:26 describes granting honour through a verb, rather than using Al Mu'izz as a name. |
| 25 | Al Mudhill | The Humiliator | [3:26](https://quran.com/3/26) | Related passage: 3:26 describes bringing low through a verb, rather than using Al Mudhill as a name. |
| 26 | As Samee' | The All-Hearing | [2:127](https://quran.com/2/127), [2:256](https://quran.com/2/256), [8:17](https://quran.com/8/17), [49:1](https://quran.com/49/1) |  |
| 27 | Al Baseer | The All-Seeing | [4:58](https://quran.com/4/58), [17:1](https://quran.com/17/1), [42:11](https://quran.com/42/11), [42:27](https://quran.com/42/27) |  |
| 28 | Al Hakam | The Judge | [6:114](https://quran.com/6/114), [22:69](https://quran.com/22/69) | Quran 6:114 uses hakaman (judge) of Allah; 22:69 describes His judging through a verb. |
| 29 | Al 'Adl | The Just | [4:40](https://quran.com/4/40), [6:115](https://quran.com/6/115) | Related passages: 4:40 affirms that Allah does not wrong anyone; 6:115 describes the justice of His word. Neither uses Al Adl as a name. |
| 30 | Al Lateef | The Subtle One | [6:103](https://quran.com/6/103), [22:63](https://quran.com/22/63), [31:16](https://quran.com/31/16), [33:34](https://quran.com/33/34) |  |
| 31 | Al Khabeer | The All-Aware | [6:18](https://quran.com/6/18), [17:30](https://quran.com/17/30), [49:13](https://quran.com/49/13), [59:18](https://quran.com/59/18) |  |
| 32 | Al Haleem | The Forbearing | [2:235](https://quran.com/2/235), [17:44](https://quran.com/17/44), [22:59](https://quran.com/22/59), [35:41](https://quran.com/35/41) |  |
| 33 | Al 'Azeem | The Magnificent | [2:255](https://quran.com/2/255), [42:4](https://quran.com/42/4), [56:96](https://quran.com/56/96) |  |
| 34 | Al Ghafoor | The Great Forgiver | [2:173](https://quran.com/2/173), [8:69](https://quran.com/8/69), [16:110](https://quran.com/16/110), [41:32](https://quran.com/41/32) |  |
| 35 | Ash Shakoor | The Appreciative | [35:30](https://quran.com/35/30), [35:34](https://quran.com/35/34), [42:23](https://quran.com/42/23), [64:17](https://quran.com/64/17) |  |
| 36 | Al 'Aliyy | The Highest | [4:34](https://quran.com/4/34), [31:30](https://quran.com/31/30), [42:4](https://quran.com/42/4), [42:51](https://quran.com/42/51) |  |
| 37 | Al Kabeer | The Greatest | [13:9](https://quran.com/13/9), [22:62](https://quran.com/22/62), [31:30](https://quran.com/31/30) |  |
| 38 | Al Hafeez | The Preserver | [11:57](https://quran.com/11/57), [34:21](https://quran.com/34/21), [42:6](https://quran.com/42/6) |  |
| 39 | Al Muqeet | The Sustainer | [4:85](https://quran.com/4/85) |  |
| 40 | Al Haseeb | The Reckoner | [4:6](https://quran.com/4/6), [4:86](https://quran.com/4/86), [33:39](https://quran.com/33/39) |  |
| 41 | Al Jaleel | The Majestic | [55:27](https://quran.com/55/27) | Related expression: 55:27 contains Dhu al-Jalal wa al-Ikram, not Al Jaleel. |
| 42 | Al Kareem | The Most Generous | [27:40](https://quran.com/27/40), [82:6](https://quran.com/82/6) |  |
| 43 | Ar Raqeeb | The Watchful One | [4:1](https://quran.com/4/1), [5:117](https://quran.com/5/117) |  |
| 44 | Al Mujeeb | The Responsive | [11:61](https://quran.com/11/61) |  |
| 45 | Al Waasi' | The All-Encompassing | [2:268](https://quran.com/2/268), [3:73](https://quran.com/3/73), [5:54](https://quran.com/5/54) |  |
| 46 | Al Hakeem | The Perfectly Wise | [31:27](https://quran.com/31/27), [46:2](https://quran.com/46/2), [57:1](https://quran.com/57/1), [66:2](https://quran.com/66/2) |  |
| 47 | Al Wadood | The Most Loving | [11:90](https://quran.com/11/90), [85:14](https://quran.com/85/14) |  |
| 48 | Al Majeed | The All-Glorious | [11:73](https://quran.com/11/73) |  |
| 49 | Al Baa'ith | The Resurrector | [22:7](https://quran.com/22/7) | Related passage: 22:7 describes raising the dead through a verb, rather than using Al Baa'ith as a name. |
| 50 | Ash Shaheed | The Witness | [4:166](https://quran.com/4/166), [22:17](https://quran.com/22/17), [41:53](https://quran.com/41/53), [48:28](https://quran.com/48/28) |  |
| 51 | Al Haqq | The Truth | [6:62](https://quran.com/6/62), [22:6](https://quran.com/22/6), [23:116](https://quran.com/23/116), [24:25](https://quran.com/24/25) |  |
| 52 | Al Wakeel | The Trustee | [3:173](https://quran.com/3/173), [4:171](https://quran.com/4/171), [28:28](https://quran.com/28/28), [73:9](https://quran.com/73/9) |  |
| 53 | Al Qawiyy | The All-Strong | [22:40](https://quran.com/22/40), [22:74](https://quran.com/22/74), [42:19](https://quran.com/42/19), [57:25](https://quran.com/57/25) |  |
| 54 | Al Mateen | The Firm | [51:58](https://quran.com/51/58) |  |
| 55 | Al Waliyy | The Protector | [4:45](https://quran.com/4/45), [7:196](https://quran.com/7/196), [42:28](https://quran.com/42/28), [45:19](https://quran.com/45/19) |  |
| 56 | Al Hameed | The Praiseworthy | [14:8](https://quran.com/14/8), [31:12](https://quran.com/31/12), [31:26](https://quran.com/31/26), [41:42](https://quran.com/41/42) |  |
| 57 | Al Muhsi | The All-Enumerating | [72:28](https://quran.com/72/28), [78:29](https://quran.com/78/29) | Related passages: these verses use verbs of counting and recording, not Al Muhsi as a name. |
| 58 | Al Mubdi | The Originator | [10:34](https://quran.com/10/34), [27:64](https://quran.com/27/64), [29:19](https://quran.com/29/19), [85:13](https://quran.com/85/13) | Related passages: these verses use verbs of originating creation, not Al Mubdi as a name. |
| 59 | Al Mu'eed | The Restorer | [10:34](https://quran.com/10/34), [27:64](https://quran.com/27/64), [29:19](https://quran.com/29/19), [85:13](https://quran.com/85/13) | Related passages: these verses use verbs of restoring creation, not Al Mu'eed as a name. |
| 60 | Al Muhyee | The Life-Giver | [7:158](https://quran.com/7/158), [15:23](https://quran.com/15/23), [30:50](https://quran.com/30/50), [57:2](https://quran.com/57/2) | Quran 30:50 uses muhyi (giver of life); the other passages express giving life through verbs. |
| 61 | Al Mumeet | The Life-Taker | [3:156](https://quran.com/3/156), [7:158](https://quran.com/7/158), [15:23](https://quran.com/15/23), [57:2](https://quran.com/57/2) | Related passages: these verses express causing death through verbs, not Al Mumeet as a name. |
| 62 | Al Haiyy | The Ever-Living | [2:255](https://quran.com/2/255), [3:2](https://quran.com/3/2), [25:58](https://quran.com/25/58), [40:65](https://quran.com/40/65) |  |
| 63 | Al Qayyoom | The Self-Subsisting | [2:255](https://quran.com/2/255), [3:2](https://quran.com/3/2), [20:111](https://quran.com/20/111) |  |
| 64 | Al Waajid | The All-Sufficient | [35:15](https://quran.com/35/15) | Related attribute: 35:15 affirms Allah's self-sufficiency using Al Ghaniyy, not Al Waajid. The explanation of Al Waajid follows al-Ghazali's Al-Maqsad al-Asna. |
| 65 | Al Maajid | The Glorious | [11:73](https://quran.com/11/73) | Related form: 11:73 contains Majeed, not Maajid. These are distinct Arabic forms with overlapping meanings. |
| 66 | Al Waahid | The Only One | [2:163](https://quran.com/2/163), [5:73](https://quran.com/5/73), [9:31](https://quran.com/9/31), [18:110](https://quran.com/18/110) |  |
| 67 | Al Ahad | The One | [112:1](https://quran.com/112/1) |  |
| 68 | As Samad | The Eternal Refuge | [112:2](https://quran.com/112/2) |  |
| 69 | Al Qaadir | The Powerful | [6:65](https://quran.com/6/65), [36:81](https://quran.com/36/81), [46:33](https://quran.com/46/33), [75:40](https://quran.com/75/40) |  |
| 70 | Al Muqtadir | The All-Powerful | [18:45](https://quran.com/18/45), [54:42](https://quran.com/54/42), [54:55](https://quran.com/54/55) |  |
| 71 | Al Muqaddim | The Expediter | [16:61](https://quran.com/16/61) | Related passage: 16:61 concerns appointed terms, not the name Al Muqaddim. The name occurs explicitly in the Arabic of Muslim 771a. |
| 72 | Al Mu’akhkhir | The Delayer | [71:4](https://quran.com/71/4) | Related passage: 71:4 describes delaying through a verb. Al Mu'akhkhir occurs explicitly in the Arabic of Muslim 771a. |
| 73 | Al Awwal | The First | [57:3](https://quran.com/57/3) |  |
| 74 | Al Aakhir | The Last | [57:3](https://quran.com/57/3) |  |
| 75 | Az Zaahir | The Manifest | [57:3](https://quran.com/57/3) |  |
| 76 | Al Baatin | The Hidden | [57:3](https://quran.com/57/3) |  |
| 77 | Al Waali | The Governor | [13:11](https://quran.com/13/11) | Quran 13:11 uses the corresponding indefinite form waalin (protector). |
| 78 | Al Muta’ali | The Supreme One | [13:9](https://quran.com/13/9) |  |
| 79 | Al Barr | The Benevolent | [52:28](https://quran.com/52/28) |  |
| 80 | At Tawwaab | Accepter of Repentance | [2:128](https://quran.com/2/128), [4:64](https://quran.com/4/64), [49:12](https://quran.com/49/12), [110:3](https://quran.com/110/3) |  |
| 81 | Al Muntaqim | The Avenger | [32:22](https://quran.com/32/22), [43:41](https://quran.com/43/41), [44:16](https://quran.com/44/16) | These passages use the plural form muntaqimun in statements of divine requital, rather than the singular Al Muntaqim. |
| 82 | Al Afuww | The Pardoner | [4:99](https://quran.com/4/99), [4:149](https://quran.com/4/149), [22:60](https://quran.com/22/60) |  |
| 83 | Ar Ra’oof | The Most Kind | [3:30](https://quran.com/3/30), [9:117](https://quran.com/9/117), [57:9](https://quran.com/57/9), [59:10](https://quran.com/59/10) |  |
| 84 | Maalik Ul Mulk | Sovereign of Dominion | [3:26](https://quran.com/3/26) |  |
| 85 | Dhu Al Jalaali Wa Al Ikraam | Majestic and Generous | [55:27](https://quran.com/55/27), [55:78](https://quran.com/55/78) | The phrase appears as dhu al-jalal in 55:27 and dhi al-jalal in 55:78 because of grammatical case. |
| 86 | Al Muqsit | The Equitable One | [7:29](https://quran.com/7/29), [3:18](https://quran.com/3/18) | Related passages: 3:18 describes upholding justice and 7:29 commands justice. Neither uses Al Muqsit as a divine name. |
| 87 | Al Jaami' | The Gatherer | [3:9](https://quran.com/3/9) |  |
| 88 | Al Ghaniyy | The Self-Sufficient | [3:97](https://quran.com/3/97), [39:7](https://quran.com/39/7), [47:38](https://quran.com/47/38), [57:24](https://quran.com/57/24) |  |
| 89 | Al Mughni | The Enricher | [9:28](https://quran.com/9/28), [53:48](https://quran.com/53/48) | Related passages: these verses describe enriching through verbs, not Al Mughni as a name. |
| 90 | Al Maani' | The Withholder | [67:21](https://quran.com/67/21) | Related passage: 67:21 describes withholding provision through a verb, not Al Maani' as a name. |
| 91 | Ad Daarr | The Distresser | [6:17](https://quran.com/6/17) | Related passage: 6:17 describes Allah's authority over harm and good, not Ad Daarr as a name. |
| 92 | An Naafi’ | The Benefactor | [10:107](https://quran.com/10/107), [6:17](https://quran.com/6/17) | Related passages: these verses describe Allah's authority over good and harm, not An Naafi' as a name. |
| 93 | An Noor | The Light | [24:35](https://quran.com/24/35) |  |
| 94 | Al Haadi | The Guide | [25:31](https://quran.com/25/31) |  |
| 95 | Al Badi' | The Incomparable Originator | [2:117](https://quran.com/2/117), [6:101](https://quran.com/6/101) |  |
| 96 | Al Baaqi | The Everlasting One | [55:27](https://quran.com/55/27) | Related passage: 55:27 expresses remaining through a verb, not Al Baaqi as a name. |
| 97 | Al Waarith | The Inheritor | [15:23](https://quran.com/15/23), [19:40](https://quran.com/19/40), [21:89](https://quran.com/21/89) | Quran 15:23 and 21:89 use plural forms of inheritor; 19:40 describes inheriting through a verb. |
| 98 | Ar Rasheed | The Rightly-Guiding | [18:10](https://quran.com/18/10) | Related passage: 18:10 asks Allah for right guidance. It does not name Him Ar Rasheed. |
| 99 | As Saboor | The Patient One | [16:61](https://quran.com/16/61), [35:45](https://quran.com/35/45) | Related passages: these verses describe delaying punishment, not Al Saboor as a name. Bukhari 6099 describes divine patience using a comparative form. |

Validation: 99 unique sequential entries and 237 individual Quran links. Every link resolves to a bundled verse; no unsupported range URL remains. Arabic strings, French fields, numbering and order are unchanged from the pre-correction corpus.
