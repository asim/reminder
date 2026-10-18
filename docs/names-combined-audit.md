# Combined names audit — 18 October 2026

Status: historical audit snapshot, recorded before correction. The approved corrections are now implemented in the names data and reference notes; see `names-sources.md` for the current source map. The original fields below are retained as audit evidence, not as the current published wording. This is not a certification of a universally agreed enumeration of 99 names.

## Scope and composition

Reviewed all 99 records together with their numbered summary files and listed Quran passages. The short label and main description come from `names/data/names.json`; the displayed summary comes from `names/data/<number>.json`. The numbered files also contain a second description which the loader does not use. Those duplicate descriptions can contradict the main record.

The React detail page displays the short label and summary, but does not display the main description. Share text and bookmark excerpts use the summary. The legacy HTML page displays both description and summary. Consequently a correction to the main description alone does not fix the modern detail page or its shared text.

The existing persisted search index may retain old wording after deployment; the PR120 review identifies that the completed-index check does not fingerprint the corpus. Previously generated hourly material is another separate output: updating source data does not rewrite historical generated content.

## Editorial rules

- Keep the Arabic identities and stable IDs; distinguish different Arabic forms even where English meanings overlap.
- Use a short label, a precise one- or two-sentence description, and a summary that expands the same meaning. Do not invent distinctions simply to make every English label unique.
- A three-word heading is a useful target, but a compound name should not lose essential meaning to satisfy a word count.
- Distinguish an actual nominal form referring to Allah from a related verb, expression, or thematic passage. An inflected or plural form should be identified explicitly.
- A verse expressing an action does not by itself establish that every derived noun is an authenticated divine name. The familiar enumeration and authentication of individual names are separate questions.
- References below are curated examples, not an exhaustive Quran concordance. Arabic-root assertions need actual linguistic support; remove unnecessary etymologies rather than retaining speculative ones.

## Sources and method

Primary Arabic verse evidence is reproduced below from the bundled `quran/data/arabic.json`; each reference links to the corresponding Quran passage for independent review. Range references were expanded for existence checks, not treated as supported application routes. Read the full verse to distinguish the speaker, referent and grammatical form.

- [Quran 112:1–4](https://quran.com/112): Ahad and Samad. [Ibn Kathir on 112:2](https://quran.com/112:2/tafsirs/en-tafisr-ibn-kathir) supports the explanation of dependence and needs.
- [Quran 34:26](https://quran.com/34/26): Fattaah in the context of judging in truth.
- [Quran 55:27](https://quran.com/55/27) and [55:78](https://quran.com/55/78): the compound Dhu al-Jalal wa al-Ikram.
- [Quran 15:23](https://quran.com/15:23), [19:40](https://quran.com/19/40), [21:89](https://quran.com/21/89): inheritance; [21:55](https://quran.com/21/55) does not support the current Waarith summary.
- [Abu Dawud 3451](https://sunnah.com/abudawud:3451): Qaabid and Baasit explicitly in Arabic.
- [Muslim 771a](https://sunnah.com/muslim:771a): Muqaddim and Mu’akhkhir explicitly in Arabic; consult the Arabic rather than relying on the English rendering of these words.
- [Muslim 2713a](https://sunnah.com/muslim:2713a): explanation of First, Last, Manifest and Hidden.
- [Bukhari 6099](https://sunnah.com/bukhari:6099): divine patience, expressed comparatively.
- [Tirmidhi 3507](https://sunnah.com/tirmidhi:3507): the familiar enumerated list; the displayed grading is weak. This does not establish all individual names as Quran occurrences.
- Existing `docs/names-sources.md` identifies al-Ghazali’s commentary for previously revised descriptions. In particular Waajid needs an explicitly attributed interpretive choice; do not turn a provisional English gloss into a claim of scholarly consensus.

## All 99 entries: findings and proposed treatment

“Nominal” means the cited passages contain the name or a corresponding nominal form referring to Allah. “Related/review” flags entries whose existing references are verbs, different expressions, or unsuitable evidence. See each note for exceptions and exact changes. “Mixed/form note” requires a per-reference distinction. These are editorial classifications, not hadith authenticity grades.

| # | Current name | Current short label | Quran evidence | Findings / proposed change |
| --- | --- | --- | --- | --- |
| 1 | Ar Rahmaan | The Beneficent | Nominal | Keep The Beneficent. Explain abundant mercy without making an exclusive theological distinction from Raheem. |
| 2 | Ar Raheem | The Merciful | Nominal | Keep The Merciful; related to Rahmaan, not a duplicate spelling. Normalize reference 3:31. |
| 3 | Al Malik | The King | Nominal | King/sovereignty is aligned; distinguish Malik from Maalik Ul Mulk by wording, without suggesting separate powers. |
| 4 | Al Quddus | The Purest | Nominal | The Holy is a clearer conventional gloss than The Purest; explain freedom from imperfection. |
| 5 | As Salaam | The Source of Peace | Nominal | Shorten to The Peace if desired; remove the unsupported claim that simply citing the name declares peace upon ourselves. |
| 6 | Al Mu'min | The inspirer of faith | Nominal | Label stresses inspiring faith, summary stresses security, main description is a long digression. Align on Giver of Security and explain confirmation of truth/faith in the fuller text. |
| 7 | Al Muhaymin | The Guardian | Nominal | Guardian/oversight is aligned; avoid inventing an exclusive distinction from Hafeez or Raqeeb. |
| 8 | Al Azeez | The Almighty | Nominal | Meaning aligned. Repair source whitespace in 4:158 and 48:7. |
| 9 | Al Jabbaar | The Compeller | Nominal | Compeller is a recognised gloss; explain irresistible authority and restoration without likening Allah to human tyranny. |
| 10 | Al Mutakabbir | The Greatest | Nominal | Prefer The Supremely Great to distinguish the wording from Kabeer; keep the existing description of greatness above imperfection. |
| 11 | Al Khaaliq | The Creator | Nominal | Replace awkward living/dead/mobile/frozen prose with creation of everything; align summary with the same point. |
| 12 | Al Baari | The Maker of Order | Nominal | The Maker is clearer than Maker of Order. Remove the claim that creation is lifted above all faults. |
| 13 | Al Musawwir | The Shaper of Beauty | Nominal | Use The Fashioner; remove the unsupported historical account about ancient sculptors and modern designers. |
| 14 | Al Ghaffaar | The Forgiving | Nominal | Include abundant/repeated forgiveness. Ghaffaar and Ghafoor overlap; do not impose a rigid theological boundary. |
| 15 | Al Qahhaar | The Subduer | Nominal | Subduer and irresistible authority are aligned. |
| 16 | Al Wahhaab | The Giver of All | Nominal | The Bestower is a concise alternative to Giver of All; retain freely giving bounty. |
| 17 | Ar Razzaaq | The Sustainer | Nominal | Prefer The Provider, matching provision in 51:58 and the revised description; distinguish wording from Muqeet without denying overlap. |
| 18 | Al Fattaah | The Opener | Nominal | Summary only describes breakthroughs. Add judging/deciding in truth, the actual context of 34:26. |
| 19 | Al 'Aleem | The Knower of all | Nominal | The All-Knowing is a concise label; knowledge of everything is aligned. |
| 20 | Al Qaabid | The Constrictor | Related/review | Quran 2:245 uses a verb, not the listed noun. Pair contraction with expansion; Abu Dawud 3451 explicitly supplies the names. |
| 21 | Al Baasit | The Reliever | Related/review | Prefer The Expander to The Reliever. 2:245 uses a verb; pair with Qaabid and cite Abu Dawud 3451 separately. |
| 22 | Al Khaafid | The Abaser | Related/review | Rewrite unclear prose about oppressors and demeanor. 95:5 is a related action, not a literal occurrence of Al Khaafid. |
| 23 | Ar Raafi' | The Exalter | Related/review | Remove The Breeder from the description. Both references describe raising through verbs, not the listed name. |
| 24 | Al Mu'iz | The Bestower of Honour | Related/review | Transliteration should reflect doubled z: Al Mu'izz. The Honourer is concise; 3:26 uses a verb. |
| 25 | Al Mudhil | The Humiliator | Related/review | Transliteration should reflect doubled l: Al Mudhill. Rewrite devalue any valiant; 3:26 uses a verb. |
| 26 | As Samee' | The Hearer of all | Nominal | Prefer The All-Hearing. Explain complete hearing without the summary's unnecessary spatial/omnipresence claim. |
| 27 | Al Baseer | The Seer of all | Nominal | Prefer The All-Seeing; seeing all things is aligned. |
| 28 | Al Hakam | The Judge | Related/review | Shorten repetitive description; add 6:114, which uses hakaman for Allah. Existing 22:69 uses the verb judges. |
| 29 | Al 'Adl | The Just | Related/review | 6:115 describes the justice of Allah's word, not a literal Al Adl. 4:40 is a clearer related passage for absence of injustice. |
| 30 | Al Lateef | The Subtle One | Nominal | Explain both subtle knowledge and gracious kindness; the current description emphasises only gentleness. |
| 31 | Al Khabeer | The All Aware | Nominal | Correct summary root k-b-r to kh-b-r (خ ب ر). Remove unrelated moralising from the main description. |
| 32 | Al Haleem | The Forbearing | Nominal | Keep The Forbearing; simplify to not hastening punishment despite knowing wrongdoing. |
| 33 | Al 'Azeem | The Magnificent | Nominal | Magnificence is aligned; avoid manufacturing a strict separation from Kabeer. |
| 34 | Al Ghafoor | The Great Forgiver | Nominal | Forgiveness is aligned; make its overlap with Ghaffaar explicit rather than inventing separate categories of sins. |
| 35 | Ash Shakoor | The Rewarder of Thankfulness | Nominal | Prefer The Appreciative; explain generous reward for sincere good rather than only rewarding thankfulness. Repair awkward description. |
| 36 | Al 'Aliyy | The Highest | Nominal | Exaltedness is aligned; keep explanation clear without adding a disputed spatial interpretation. |
| 37 | Al Kabeer | The Greatest | Nominal | The Greatest is aligned with the revised description; meaningful overlap with other greatness names is expected. |
| 38 | Al Hafeez | The Preserver | Nominal | Clarify the vague preserves His good wording: preserving/guarding and knowledge of deeds. |
| 39 | Al Muqeet | The Sustainer | Nominal | The Sustainer fits the existing revised description. Correct incomplete q-t root to q-w-t (ق و ت); acknowledge interpretations involving power/watchfulness in 4:85. |
| 40 | Al Haseeb | The Reckoner, The Sufficient | Nominal | Use The Reckoner as the short label and retain sufficiency in the fuller explanation; avoid a multi-label heading. |
| 41 | Al Jaleel | The Majestic | Related/review | 55:27 contains Dhu al-Jalal, a related expression, not Al Jaleel. Remove 39:14 as unrelated; 7:143 is thematic, not a name occurrence. |
| 42 | Al Kareem | The Most Generous | Nominal | Generosity is aligned; keep short label and explanation consistent. |
| 43 | Ar Raqeeb | The Watchful One | Nominal | Watchfulness is aligned; overlap with Shaheed/Hafeez is legitimate. |
| 44 | Al Mujeeb | The Responder to Prayer | Nominal | Trim trailing whitespace in transliteration. Correct incomplete J-B root to j-w-b (ج و ب). The Responsive is concise. |
| 45 | Al Waasi' | The All-Encompassing | Nominal | Remove fills every void/contains all spatial claims and omnipresence wording. Explain encompassing knowledge, mercy and bounty. |
| 46 | Al Hakeem | The Perfectly Wise | Nominal | Replace repeated dazzled prose with wisdom in judgment, action and decree. |
| 47 | Al Wadood | The Most Loving | Nominal | Replace remarks about translation difficulty with a direct explanation of divine love; preserve corrected Wadood spelling. |
| 48 | Al Majeed | The Glorious One, The Most Honorable | Nominal | Use The All-Glorious or The Glorious; the current heading is six words. Do not invent an absolute distinction from Maajid. |
| 49 | Al Baa'ith | The Resurrector | Related/review | 22:7 describes resurrection through a verb; label it a related attribute, not a literal occurrence of Al Baa'ith. |
| 50 | Ash Shaheed | The Witness | Nominal | Witness is aligned; remove the unnecessary omnipresence claim from the explanation. |
| 51 | Al Haqq | The Truth | Nominal | Truth/reality is aligned. |
| 52 | Al Wakeel | The Trustee | Nominal | Trustee/reliance is aligned; do not imply that reliance removes human responsibility. |
| 53 | Al Qawiyy | The Possessor of all strength | Nominal | The All-Strong is concise; explain complete strength without bodily comparisons. |
| 54 | Al Mateen | The Forceful | Nominal | The Firm is clearer than The Forceful; strength is unfailing, not merely forceful conduct. |
| 55 | Al Waliyy | The Protector | Nominal | Replace rambling main description with protecting and supporting believers. Waliyy and Waali are distinct forms with overlapping meanings. |
| 56 | Al Hameed | The Praised | Nominal | Praise is aligned; The Praiseworthy is more idiomatic than The Praised if changing the label. |
| 57 | Al Muhsi | The All-Enumerating | Related/review | Correct root حصر to ح ص ي. 72:28 and 78:29 use verbs of enumeration. 82:10–12 describes recording angels and also creates an unsupported range URL. |
| 58 | Al Mubdi | The Originator | Related/review | Origination is aligned. All existing references use verbs, not the listed noun. |
| 59 | Al Mu'eed | The Restorer | Related/review | Correct the description's direction: restoring creation/resurrection, not bringing the living back to death. Remove the summary's unsupported focus on second chances/forgiveness. References use verbs. |
| 60 | Al Muhyee | The Giver of life | Mixed/form note | 30:50 has the nominal form muhyi; other references use verbs. Remove unnecessary biological and resurrection details from the description. |
| 61 | Al Mumeet | The Taker of life | Related/review | Life-taking is aligned; all listed references use verbs rather than the noun Al Mumeet. |
| 62 | Al Haiyy | The Ever Living | Nominal | Remove the long unrelated quotation about disbelief from the description. Explain life without death, sleep or dependence. |
| 63 | Al Qayyoom | The Self Existing | Nominal | Description must include sustaining all existence as well as self-subsistence; summary already includes this. The Sustainer is concise but may legitimately overlap other labels. |
| 64 | Al Waajid | The Perceiver | Related/review | Substantive mismatch: Perceiver/summary of perception versus description of lacking nothing. Resolve using a named commentary before choosing a label. 38:44 says Allah found Ayyub patient; it does not call Allah Al Waajid. |
| 65 | Al Maajid | The Glorious | Related/review | Glory is aligned; 11:73 contains Majeed, not Maajid. Mark as a related form and retain the distinction in Arabic without inventing a rigid English distinction. |
| 66 | Al Waahid | The Only One | Nominal | Oneness is aligned; Waahid and Ahad are different words expressing overlapping truths, not accidental duplicates. |
| 67 | Al Ahad | The One | Nominal | Remove the false covenant etymology: أحد (Ahad) is not عهد (ahd). Rebuild summary around unique oneness and 112:1–4. |
| 68 | As Samad | The Eternal Refuge | Nominal | Keep agreed The Eternal Refuge and revised main description. Update separate summary to state both needing nothing and all creation turning to Him in need; current summary is too generic. |
| 69 | Al Qaadir | The Powerful | Nominal | Power is aligned; preserve meaningful overlap with Muqtadir. |
| 70 | Al Muqtadir | The Creator of all power | Nominal | Prefer The All-Powerful over Creator of all power; describe complete power and capability. |
| 71 | Al Muqaddim | The Expediter | Related/review | Remove unrelated 17:34. 16:61 concerns appointed terms, not literal Al Muqaddim. Muslim 771a contains the name explicitly in Arabic. |
| 72 | Al Mu’akhkhir | The Delayer | Related/review | 71:4 uses a verb of delaying, not the listed noun. Use Muslim 771a for the explicit name; keep corrected transliteration. |
| 73 | Al Awwal | The First | Nominal | 57:3 is direct. Use Muslim 2713a's explanation that nothing precedes Him. |
| 74 | Al Aakhir | The Last | Nominal | 57:3 is direct. Use Muslim 2713a's explanation that nothing comes after Him. |
| 75 | Az Zaahir | The Manifest | Nominal | 57:3 is direct. Explain with Muslim 2713a; avoid silently choosing a contested spatial paraphrase. |
| 76 | Al Baatin | The Hidden | Nominal | 57:3 is direct. Explain with Muslim 2713a; do not equate hiddenness with being absent or unknown in every sense. |
| 77 | Al Waali | The Governor | Mixed/form note | 13:11 has waalin, a corresponding indefinite form meaning protector. Remove unrelated 22:7, which is about resurrection. |
| 78 | Al Muta’ali | The Supreme One | Nominal | 13:9 directly contains al-muta'al. Supremacy/exaltedness is aligned. |
| 79 | Al Barr | The Doer of Good | Nominal | The Benevolent is concise. Remove duplicated non-Muslims text and repair the main description. |
| 80 | At Tawwaab | The Guide to Repentance | Nominal | Change Guide to Repentance to Accepter of Repentance; align summary with both enabling repentance and accepting it, as in the revised description. |
| 81 | Al Muntaqim | The Avenger | Mixed/form note | Expand the three-word description to just requital for wrongdoing. References use plural muntaqimun, not singular Al Muntaqim; show that grammatical distinction. |
| 82 | Al Afuww | The Forgiver | Nominal | Change The Forgiver to The Pardoner, matching the revised description of effacing sins; forgiveness and pardon still overlap. |
| 83 | Ar Ra’oof | The Most Kind | Nominal | Kindness/compassion is aligned. |
| 84 | Maalik Ul Mulk | Master of the Kingdom | Nominal | Sovereign of Dominion is a concise option; 3:26 contains the phrase directly. |
| 85 | Dhu Al Jalaali Wa Al Ikraam | Possessor of Majesty and Bounty | Nominal | Remove God's dual nature. Name is explicitly in 55:27 and 55:78, not merely hadith. Explain majesty and generosity/honour. Do not force a three-word label if it loses half the compound meaning. |
| 86 | Al Muqsit | The Equitable One | Related/review | 3:18 describes upholding justice; 7:29 commands justice. Neither literally calls Allah Al Muqsit. Mark as related attributes. |
| 87 | Al Jaami' | The Gatherer | Nominal | 3:9 has jami', a direct nominal form referring to gathering people. Meaning is aligned. |
| 88 | Al Ghaniyy | The Rich One | Nominal | Prefer The Self-Sufficient to The Rich One, avoiding a merely financial reading. Explain no dependence on creation. |
| 89 | Al Mughni | The Enricher | Related/review | Replace separate description Inflictor of Help/Grants Victory and summary about triumph. Enriching/giving sufficiency is the relevant meaning. 9:28 and suggested 53:48 use verbs. |
| 90 | Al Maani' | The Preventer of harm | Related/review | Change Preventer of harm to The Withholder. 67:21 concerns withholding provision; it does not promise universal protection from harm or contain this name as a noun. |
| 91 | Ad Daaarr | The Creator of the harmful | Related/review | Correct Daaarr to Daarr and false root Daal-Ra-Ayn to ض ر ر. Present hardship in relation to benefit and divine wisdom, with no implication of injustice. 6:17 is related, not a literal name occurrence. |
| 92 | An Naafi’ | The Bestower of Benefits | Related/review | 30:37 concerns provision rather than a literal name. Prefer 10:107 as a related passage on harm/good, alongside 6:17; explain benefit, paired with Daarr. |
| 93 | An Noor | The Light | Nominal | 24:35 directly supplies Light. Keep the explanation grounded in the verse and distinguish any interpretive imagery from the text itself. |
| 94 | Al Haadi | The Guider | Nominal | 25:31 has haadiyan, a direct nominal form. The Guide is more idiomatic than The Guider. |
| 95 | Al Badi' | The Incomparable Originator | Nominal | 2:117 and 6:101 contain badi' directly. Explain origination without a prior model; overlap with Mubdi is not a data error. |
| 96 | Al Baaqi | The Everlasting One | Related/review | Remove the unqualified claim that all creation ceases to exist. 55:26–27 contrasts earthly mortality with Allah's remaining; 55:27 uses a verb, not Al Baaqi. |
| 97 | Al Waarith | The Inheritor | Mixed/form note | Correct the summary's false 21:55 attribution. Existing 15:23 has plural al-warithun; 19:40 and 21:89 provide additional relevant passages. |
| 98 | Ar Rasheed | The Most Righteous Guide | Related/review | 2:256 mentions right guidance, not Allah as Ar Rasheed. Correct the summary's claim that راشد is a root: root is ر ش د. Avoid Rightly Guided wording that suggests Allah receives guidance. |
| 99 | As Saboor | The Patient One | Related/review | Current references command/commend human patience, not a literal divine name. Use 16:61/35:45 as related passages about delaying punishment; Bukhari 6099 supports divine patience, using a comparative form rather than Al Saboor. |

## Composed records and verse evidence

The following is the audited snapshot, including the otherwise unused numbered-file description. It makes disagreements between layers directly reviewable. Proposed corrections are in the table above; the quoted application text below is not endorsed as accurate.

### 1. Ar Rahmaan — الرَّحْمَنُ

**Short label:** The Beneficent

**Main description:** He who wills goodness and mercy for all His creatures

**Numbered-file description (unused):** Ar Rahmaan is one of the 99 names of Allah (Asma'ul Husna) in Islam, meaning 'The Most Gracious' or 'The Beneficent'. This name highlights Allah's boundless mercy and kindness towards all of creation, reflecting His nurturing and compassionate nature.

**Displayed summary:** The name Ar Rahmaan emphasizes Allah's absolute mercy and generosity. It signifies that His grace encompasses all beings, providing sustenance and care without discrimination. The name originates from the Arabic root word 'rahma', which denotes mercy and compassion. Ar Rahmaan serves as a reminder of Allah's infinite kindness and serves an important role in Islamic belief, where mercy is seen as an essential attribute of God.

**Audit:** Keep The Beneficent. Explain abundant mercy without making an exclusive theological distinction from Raheem.

- [1:3](https://quran.com/1/3): ٱلرَّحۡمَٰنِ ٱلرَّحِيمِ
- [17:110](https://quran.com/17/110): قُلِ ٱدۡعُواْ ٱللَّهَ أَوِ ٱدۡعُواْ ٱلرَّحۡمَٰنَۖ أَيّٗا مَّا تَدۡعُواْ فَلَهُ ٱلۡأَسۡمَآءُ ٱلۡحُسۡنَىٰۚ وَلَا تَجۡهَرۡ بِصَلَاتِكَ وَلَا تُخَافِتۡ بِهَا وَٱبۡتَغِ بَيۡنَ ذَٰلِكَ سَبِيلٗا
- [55:1](https://quran.com/55/1): ٱلرَّحۡمَٰنُ

### 2. Ar Raheem — الرَّحِيمُ

**Short label:** The Merciful

**Main description:** He who acts with extreme kindness

**Numbered-file description (unused):** Ar-Raheem is one of the names of Allah in Islam, which translates to 'The Most Merciful' or 'The Beneficent'. It emphasizes Allah's compassion and mercy towards His creations.

**Displayed summary:** The name Ar-Raheem originates from the Arabic root 'R-H-M', which denotes mercy and compassion. It signifies Allah's attribute of bestowing kindness and forgiveness to His creatures, highlighting the importance of mercy in the Islamic faith. Ar-Raheem often appears alongside Ar-Rahman, which also denotes mercy but encompasses a broader, universal aspect, while Ar-Raheem focuses more on a specific, heartfelt mercy towards the believers.

**Audit:** Keep The Merciful; related to Rahmaan, not a duplicate spelling. Normalize reference 3:31.

- [2:163](https://quran.com/2/163): وَإِلَٰهُكُمۡ إِلَٰهٞ وَٰحِدٞۖ لَّآ إِلَٰهَ إِلَّا هُوَ ٱلرَّحۡمَٰنُ ٱلرَّحِيمُ
- [3:31](https://quran.com/3/31): قُلۡ إِن كُنتُمۡ تُحِبُّونَ ٱللَّهَ فَٱتَّبِعُونِي يُحۡبِبۡكُمُ ٱللَّهُ وَيَغۡفِرۡ لَكُمۡ ذُنُوبَكُمۡۚ وَٱللَّهُ غَفُورٞ رَّحِيمٞ
- [4:100](https://quran.com/4/100): ۞ وَمَن يُهَاجِرۡ فِي سَبِيلِ ٱللَّهِ يَجِدۡ فِي ٱلۡأَرۡضِ مُرَٰغَمٗا كَثِيرٗا وَسَعَةٗۚ وَمَن يَخۡرُجۡ مِنۢ بَيۡتِهِۦ مُهَاجِرًا إِلَى ٱللَّهِ وَرَسُولِهِۦ ثُمَّ يُدۡرِكۡهُ ٱلۡمَوۡتُ فَقَدۡ وَقَعَ أَجۡرُهُۥ عَلَى ٱللَّهِۗ وَكَانَ ٱللَّهُ غَفُورٗا رَّحِيمٗا
- [5:3](https://quran.com/5/3): حُرِّمَتۡ عَلَيۡكُمُ ٱلۡمَيۡتَةُ وَٱلدَّمُ وَلَحۡمُ ٱلۡخِنزِيرِ وَمَآ أُهِلَّ لِغَيۡرِ ٱللَّهِ بِهِۦ وَٱلۡمُنۡخَنِقَةُ وَٱلۡمَوۡقُوذَةُ وَٱلۡمُتَرَدِّيَةُ وَٱلنَّطِيحَةُ وَمَآ أَكَلَ ٱلسَّبُعُ إِلَّا مَا ذَكَّيۡتُمۡ وَمَا ذُبِحَ عَلَى ٱلنُّصُبِ وَأَن تَسۡتَقۡسِمُواْ بِٱلۡأَزۡلَٰمِۚ ذَٰلِكُمۡ فِسۡقٌۗ ٱلۡيَوۡمَ يَئِسَ ٱلَّذِينَ كَفَرُواْ مِن دِينِكُمۡ فَلَا تَخۡشَوۡهُمۡ وَٱخۡشَوۡنِۚ ٱلۡيَوۡمَ أَكۡمَلۡتُ لَكُمۡ دِينَكُمۡ وَأَتۡمَمۡتُ عَلَيۡكُمۡ نِعۡمَتِي وَرَضِيتُ لَكُمُ ٱلۡإِسۡلَٰمَ دِينٗاۚ فَمَنِ ٱضۡطُرَّ فِي مَخۡمَصَةٍ غَيۡرَ مُتَجَانِفٖ لِّإِثۡمٖ فَإِنَّ ٱللَّهَ غَفُورٞ رَّحِيمٞ

### 3. Al Malik — الْمَلِكُ

**Short label:** The King

**Main description:** Al-Malik translates to The King or The Sovereign. This name of Allah signifies His absolute dominion and ownership over all of creation, highlighting His supreme authority.

**Numbered-file description (unused):** Al Malik translates to 'The King' in Arabic. It is one of the names of Allah in Islam, denoting His sovereignty and authority over all creation.

**Displayed summary:** Al Malik emphasizes God's ultimate dominion and control over the universe. This name reflects His majesty, power, and the structure He has established in creation. It originates from the Quran, where God's sovereignty is highlighted, demonstrating that He is the rightful ruler and king over all that exists.

**Audit:** King/sovereignty is aligned; distinguish Malik from Maalik Ul Mulk by wording, without suggesting separate powers.

- [20:114](https://quran.com/20/114): فَتَعَٰلَى ٱللَّهُ ٱلۡمَلِكُ ٱلۡحَقُّۗ وَلَا تَعۡجَلۡ بِٱلۡقُرۡءَانِ مِن قَبۡلِ أَن يُقۡضَىٰٓ إِلَيۡكَ وَحۡيُهُۥۖ وَقُل رَّبِّ زِدۡنِي عِلۡمٗا
- [23:116](https://quran.com/23/116): فَتَعَٰلَى ٱللَّهُ ٱلۡمَلِكُ ٱلۡحَقُّۖ لَآ إِلَٰهَ إِلَّا هُوَ رَبُّ ٱلۡعَرۡشِ ٱلۡكَرِيمِ
- [59:23](https://quran.com/59/23): هُوَ ٱللَّهُ ٱلَّذِي لَآ إِلَٰهَ إِلَّا هُوَ ٱلۡمَلِكُ ٱلۡقُدُّوسُ ٱلسَّلَٰمُ ٱلۡمُؤۡمِنُ ٱلۡمُهَيۡمِنُ ٱلۡعَزِيزُ ٱلۡجَبَّارُ ٱلۡمُتَكَبِّرُۚ سُبۡحَٰنَ ٱللَّهِ عَمَّا يُشۡرِكُونَ
- [62:1](https://quran.com/62/1): يُسَبِّحُ لِلَّهِ مَا فِي ٱلسَّمَٰوَٰتِ وَمَا فِي ٱلۡأَرۡضِ ٱلۡمَلِكِ ٱلۡقُدُّوسِ ٱلۡعَزِيزِ ٱلۡحَكِيمِ

### 4. Al Quddus — الْقُدُّوسُ

**Short label:** The Purest

**Main description:** The One who is pure from any imperfection and clear from children and adversaries

**Numbered-file description (unused):** Al Quddus, meaning 'The Most Sacred' or 'The Most Holy', is one of the names (Asma’ul Husna) of Allah in Islam.

**Displayed summary:** The name Al Quddus emphasizes the purity and sanctity of Allah. It signifies that Allah is free from any imperfection and holds perfect holiness. This name is often invoked in prayers and supplications to highlight the divine essence of purity and to seek closeness to Allah through recognition of His holiness.

**Audit:** The Holy is a clearer conventional gloss than The Purest; explain freedom from imperfection.

- [59:23](https://quran.com/59/23): هُوَ ٱللَّهُ ٱلَّذِي لَآ إِلَٰهَ إِلَّا هُوَ ٱلۡمَلِكُ ٱلۡقُدُّوسُ ٱلسَّلَٰمُ ٱلۡمُؤۡمِنُ ٱلۡمُهَيۡمِنُ ٱلۡعَزِيزُ ٱلۡجَبَّارُ ٱلۡمُتَكَبِّرُۚ سُبۡحَٰنَ ٱللَّهِ عَمَّا يُشۡرِكُونَ
- [62:1](https://quran.com/62/1): يُسَبِّحُ لِلَّهِ مَا فِي ٱلسَّمَٰوَٰتِ وَمَا فِي ٱلۡأَرۡضِ ٱلۡمَلِكِ ٱلۡقُدُّوسِ ٱلۡعَزِيزِ ٱلۡحَكِيمِ

### 5. As Salaam — السَّلاَمُ

**Short label:** The Source of Peace

**Main description:** The one who provides peace, calm and serenity. By citing this name, we declare peace for ourselves and for those around us.

**Numbered-file description (unused):** As Salaam is one of the 99 names of Allah in Islam, meaning 'The Source of Peace' or 'The Flawless'. It signifies that Allah is the ultimate giver of peace, safety, and tranquility.

**Displayed summary:** Originating from Arabic, the name 'As Salaam' highlights Allah's nature as the one who provides peace and serenity to His creation. It emphasizes a key attribute of Allah, fostering a sense of calm and security for believers. The invocation of this name is often used in prayers and supplications, seeking divine peace in one's life.

**Audit:** Shorten to The Peace if desired; remove the unsupported claim that simply citing the name declares peace upon ourselves.

- [59:23](https://quran.com/59/23): هُوَ ٱللَّهُ ٱلَّذِي لَآ إِلَٰهَ إِلَّا هُوَ ٱلۡمَلِكُ ٱلۡقُدُّوسُ ٱلسَّلَٰمُ ٱلۡمُؤۡمِنُ ٱلۡمُهَيۡمِنُ ٱلۡعَزِيزُ ٱلۡجَبَّارُ ٱلۡمُتَكَبِّرُۚ سُبۡحَٰنَ ٱللَّهِ عَمَّا يُشۡرِكُونَ

### 6. Al Mu'min — الْمُؤْمِنُ

**Short label:** The inspirer of faith

**Main description:** This name is mysterious even for Muslims because \"Al-Mu'min\" comes from the name \"Al-iman\" which means to believe, from where comes the name \"Al-Mu'minoune\" which means believers, it is is the name of the Koran for people who believe in 6 things: the oneness of God, the existence of angels, holy books (Torah, Gospel and Psalms (Zabor) unmodified and the Koran), the prophets, the last day (the day of judgment) and fate its evil and its good. The exact meaning of this name is very difficult to pin down in its original meaning but can be understood from a deeper and more mystical point of view.

**Numbered-file description (unused):** Al Mu'min is one of the Names of Allah in Islam, meaning 'The Granter of Security' or 'The Infuser of Faith'.

**Displayed summary:** Al Mu'min signifies Allah's role in providing peace, safety, and faith to His creatures. It emphasizes His ability to grant security not only in a worldly sense but also in the spiritual realm, instilling belief within hearts. This name is derived from the Arabic root word 'aman', which means security and faith. It reflects Allah's mercy and protection over His creations.

**Audit:** Label stresses inspiring faith, summary stresses security, main description is a long digression. Align on Giver of Security and explain confirmation of truth/faith in the fuller text.

- [59:23](https://quran.com/59/23): هُوَ ٱللَّهُ ٱلَّذِي لَآ إِلَٰهَ إِلَّا هُوَ ٱلۡمَلِكُ ٱلۡقُدُّوسُ ٱلسَّلَٰمُ ٱلۡمُؤۡمِنُ ٱلۡمُهَيۡمِنُ ٱلۡعَزِيزُ ٱلۡجَبَّارُ ٱلۡمُتَكَبِّرُۚ سُبۡحَٰنَ ٱللَّهِ عَمَّا يُشۡرِكُونَ

### 7. Al Muhaymin — الْمُهَيْمِنُ

**Short label:** The Guardian

**Main description:** Whoever dominates His creation and by His omnipotence in all aspects, no state or aspect in this world or the other escapes His will.

**Numbered-file description (unused):** Al Muhaymin is one of the names of Allah in Islam, which translates to 'The Guardian' or 'The Preserver of Safety'. It signifies God's oversight and protection over His creation.

**Displayed summary:** The name Al Muhaymin originates from the Arabic root 'h-m-n', which means to guard, protect, or supervise. In the context of Islamic theology, it reflects Allah's ability to safeguard and preserve the universe and everything within it. It embodies reassurance for believers that their affairs are under divine protection and that Allah is aware of all actions and intentions.

**Audit:** Guardian/oversight is aligned; avoid inventing an exclusive distinction from Hafeez or Raqeeb.

- [59:23](https://quran.com/59/23): هُوَ ٱللَّهُ ٱلَّذِي لَآ إِلَٰهَ إِلَّا هُوَ ٱلۡمَلِكُ ٱلۡقُدُّوسُ ٱلسَّلَٰمُ ٱلۡمُؤۡمِنُ ٱلۡمُهَيۡمِنُ ٱلۡعَزِيزُ ٱلۡجَبَّارُ ٱلۡمُتَكَبِّرُۚ سُبۡحَٰنَ ٱللَّهِ عَمَّا يُشۡرِكُونَ

### 8. Al Azeez — الْعَزِيزُ

**Short label:** The Almighty

**Main description:** This name signifies Allah's absolute strength, might, and dominion, emphasizing that He is the one who cannot be overcome.

**Numbered-file description (unused):** Al Azeez translates to 'The Mighty' or 'The Invincible' in Arabic. It signifies one who has strength, power, and authority, and is often used to denote Allah's supreme power over all creation.

**Displayed summary:** Al Azeez is one of the 99 names of Allah in Islam, reflecting His infinite might and strength. Originating from the Arabic root word 'عز' (ʿizz), which means honor, glory, and strength, this name emphasizes the greatness of Allah and His ability to overpower all opposition, granting dignity and protection to His creation.

**Audit:** Meaning aligned. Repair source whitespace in 4:158 and 48:7.

- [3:6](https://quran.com/3/6): هُوَ ٱلَّذِي يُصَوِّرُكُمۡ فِي ٱلۡأَرۡحَامِ كَيۡفَ يَشَآءُۚ لَآ إِلَٰهَ إِلَّا هُوَ ٱلۡعَزِيزُ ٱلۡحَكِيمُ
- [4:158](https://quran.com/4/158): بَل رَّفَعَهُ ٱللَّهُ إِلَيۡهِۚ وَكَانَ ٱللَّهُ عَزِيزًا حَكِيمٗا
- [9:40](https://quran.com/9/40): إِلَّا تَنصُرُوهُ فَقَدۡ نَصَرَهُ ٱللَّهُ إِذۡ أَخۡرَجَهُ ٱلَّذِينَ كَفَرُواْ ثَانِيَ ٱثۡنَيۡنِ إِذۡ هُمَا فِي ٱلۡغَارِ إِذۡ يَقُولُ لِصَٰحِبِهِۦ لَا تَحۡزَنۡ إِنَّ ٱللَّهَ مَعَنَاۖ فَأَنزَلَ ٱللَّهُ سَكِينَتَهُۥ عَلَيۡهِ وَأَيَّدَهُۥ بِجُنُودٖ لَّمۡ تَرَوۡهَا وَجَعَلَ كَلِمَةَ ٱلَّذِينَ كَفَرُواْ ٱلسُّفۡلَىٰۗ وَكَلِمَةُ ٱللَّهِ هِيَ ٱلۡعُلۡيَاۗ وَٱللَّهُ عَزِيزٌ حَكِيمٌ
- [48:7](https://quran.com/48/7): وَلِلَّهِ جُنُودُ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِۚ وَكَانَ ٱللَّهُ عَزِيزًا حَكِيمًا

### 9. Al Jabbaar — الْجَبَّارُ

**Short label:** The Compeller

**Main description:** The grandiose, the ultimate powerful by His greatness and His unlimited strength. When He destroys the oppressors (unjust and sinners) He applies His sentence with righteousness.

**Numbered-file description (unused):** Al Jabbaar is one of the names of Allah in Islam, derived from the Arabic root 'j-b-r', which means to compel, to repair, or to make whole. It signifies Allah's attribute of power and authority, to enact His will and establish what is right.

**Displayed summary:** Al Jabbaar, meaning 'The Compeller' or 'The Restorer', reflects Allah's supreme ability to enforce His will upon creation. It emphasizes His power to rectify and restore all that is broken or unjust. This name highlights Allah's dominion over the universe, ensuring that justice and order prevail according to His divine plan.

**Audit:** Compeller is a recognised gloss; explain irresistible authority and restoration without likening Allah to human tyranny.

- [59:23](https://quran.com/59/23): هُوَ ٱللَّهُ ٱلَّذِي لَآ إِلَٰهَ إِلَّا هُوَ ٱلۡمَلِكُ ٱلۡقُدُّوسُ ٱلسَّلَٰمُ ٱلۡمُؤۡمِنُ ٱلۡمُهَيۡمِنُ ٱلۡعَزِيزُ ٱلۡجَبَّارُ ٱلۡمُتَكَبِّرُۚ سُبۡحَٰنَ ٱللَّهِ عَمَّا يُشۡرِكُونَ

### 10. Al Mutakabbir — الْمُتَكَبِّرُ

**Short label:** The Greatest

**Main description:** Allah alone possesses absolute majesty and greatness; no creature can claim that greatness for itself.

**Numbered-file description (unused):** Al Mutakabbir is one of the names of Allah in Islam, derived from the Arabic root 'k-b-r' which means 'to be great' or 'to be proud'.

**Displayed summary:** Al Mutakabbir translates to 'The Supreme' or 'The Majestic'. It signifies Allah's greatness and His absolute sovereignty over all created beings, emphasizing His power and majesty. This name encourages humility in humans, reminding them of their reliance on Allah's greatness.

**Audit:** Prefer The Supremely Great to distinguish the wording from Kabeer; keep the existing description of greatness above imperfection.

- [59:23](https://quran.com/59/23): هُوَ ٱللَّهُ ٱلَّذِي لَآ إِلَٰهَ إِلَّا هُوَ ٱلۡمَلِكُ ٱلۡقُدُّوسُ ٱلسَّلَٰمُ ٱلۡمُؤۡمِنُ ٱلۡمُهَيۡمِنُ ٱلۡعَزِيزُ ٱلۡجَبَّارُ ٱلۡمُتَكَبِّرُۚ سُبۡحَٰنَ ٱللَّهِ عَمَّا يُشۡرِكُونَ

### 11. Al Khaaliq — الْخَالِقُ

**Short label:** The Creator

**Main description:** The Creator of all existences, living or dead, mobile or frozen.

**Numbered-file description (unused):** Al Khaaliq is one of the names of Allah in Islam, meaning 'The Creator' or 'The Maker'. It signifies the attribute of God as the one who creates and brings forth everything that exists.

**Displayed summary:** The name Al Khaaliq originates from the Arabic root 'Khalq', which means creation. In Islamic theology, it emphasizes God's role as the ultimate creator of the universe, reflecting His divine power and wisdom in bringing things into existence. This name illustrates the belief that all creation is a manifestation of God's will and creativity.

**Audit:** Replace awkward living/dead/mobile/frozen prose with creation of everything; align summary with the same point.

- [6:102](https://quran.com/6/102): ذَٰلِكُمُ ٱللَّهُ رَبُّكُمۡۖ لَآ إِلَٰهَ إِلَّا هُوَۖ خَٰلِقُ كُلِّ شَيۡءٖ فَٱعۡبُدُوهُۚ وَهُوَ عَلَىٰ كُلِّ شَيۡءٖ وَكِيلٞ
- [13:16](https://quran.com/13/16): قُلۡ مَن رَّبُّ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِ قُلِ ٱللَّهُۚ قُلۡ أَفَٱتَّخَذۡتُم مِّن دُونِهِۦٓ أَوۡلِيَآءَ لَا يَمۡلِكُونَ لِأَنفُسِهِمۡ نَفۡعٗا وَلَا ضَرّٗاۚ قُلۡ هَلۡ يَسۡتَوِي ٱلۡأَعۡمَىٰ وَٱلۡبَصِيرُ أَمۡ هَلۡ تَسۡتَوِي ٱلظُّلُمَٰتُ وَٱلنُّورُۗ أَمۡ جَعَلُواْ لِلَّهِ شُرَكَآءَ خَلَقُواْ كَخَلۡقِهِۦ فَتَشَٰبَهَ ٱلۡخَلۡقُ عَلَيۡهِمۡۚ قُلِ ٱللَّهُ خَٰلِقُ كُلِّ شَيۡءٖ وَهُوَ ٱلۡوَٰحِدُ ٱلۡقَهَّٰرُ
- [39:62](https://quran.com/39/62): ٱللَّهُ خَٰلِقُ كُلِّ شَيۡءٖۖ وَهُوَ عَلَىٰ كُلِّ شَيۡءٖ وَكِيلٞ
- [40:62](https://quran.com/40/62): ذَٰلِكُمُ ٱللَّهُ رَبُّكُمۡ خَٰلِقُ كُلِّ شَيۡءٖ لَّآ إِلَٰهَ إِلَّا هُوَۖ فَأَنَّىٰ تُؤۡفَكُونَ
- [59:24](https://quran.com/59/24): هُوَ ٱللَّهُ ٱلۡخَٰلِقُ ٱلۡبَارِئُ ٱلۡمُصَوِّرُۖ لَهُ ٱلۡأَسۡمَآءُ ٱلۡحُسۡنَىٰۚ يُسَبِّحُ لَهُۥ مَا فِي ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِۖ وَهُوَ ٱلۡعَزِيزُ ٱلۡحَكِيمُ

### 12. Al Baari — الْبَارِئُ

**Short label:** The Maker of Order

**Main description:** He who lifts His creation above all faults or imperfections.

**Numbered-file description (unused):** Al Baari is one of the names of Allah in Islamic tradition, meaning 'The Evolver' or 'The Creator'. It signifies Allah's role in creating all things in their perfect forms and stages.

**Displayed summary:** The name Al Baari derives from the Arabic root 'ب ر أ' (B-R-ʾ), which conveys the concept of bringing into existence or creating. It emphasizes Allah's authority and power in shaping the universe and everything within it from non-existence to existence, tailoring each creation with purpose and precision.

**Audit:** The Maker is clearer than Maker of Order. Remove the claim that creation is lifted above all faults.

- [59:24](https://quran.com/59/24): هُوَ ٱللَّهُ ٱلۡخَٰلِقُ ٱلۡبَارِئُ ٱلۡمُصَوِّرُۖ لَهُ ٱلۡأَسۡمَآءُ ٱلۡحُسۡنَىٰۚ يُسَبِّحُ لَهُۥ مَا فِي ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِۖ وَهُوَ ٱلۡعَزِيزُ ٱلۡحَكِيمُ

### 13. Al Musawwir — الْمُصَوِّرُ

**Short label:** The Shaper of Beauty

**Main description:** The one who gives the look and the appearance to each creation from the inside as well as from the outside. In lost Arabic, literally, this word means sculptor and, in modern Arabic, literally means designer.

**Numbered-file description (unused):** Al Musawwir translates to 'The Fashioner' or 'The Shaper'. It refers to Allah as the one who shapes and forms creation according to His will.

**Displayed summary:** Al Musawwir is one of the names of Allah in Islam, highlighting His role in shaping and defining the forms of all beings. The name emphasizes the intricate design and beauty of creation, illustrating that all forms and features in the universe are crafted by the Creator. The understanding of this name encourages believers to recognize the deliberate and purposeful creation of the world around them.

**Audit:** Use The Fashioner; remove the unsupported historical account about ancient sculptors and modern designers.

- [59:24](https://quran.com/59/24): هُوَ ٱللَّهُ ٱلۡخَٰلِقُ ٱلۡبَارِئُ ٱلۡمُصَوِّرُۖ لَهُ ٱلۡأَسۡمَآءُ ٱلۡحُسۡنَىٰۚ يُسَبِّحُ لَهُۥ مَا فِي ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِۖ وَهُوَ ٱلۡعَزِيزُ ٱلۡحَكِيمُ

### 14. Al Ghaffaar — الْغَفَّارُ

**Short label:** The Forgiving

**Main description:** Is the One who forgives what He wants to whom He wants.

**Numbered-file description (unused):** Al Ghaffaar is one of the names of Allah in Islam, derived from the Arabic root ' غ ف ر ' (ghafara), which means to forgive or to conceal.

**Displayed summary:** The meaning of Al Ghaffaar is 'The Constant Forgiver' or 'The Forgiving One'. It emphasizes Allah's willingness to forgive sins repeatedly and His mercy towards humanity. This name signifies that Allah continually provides opportunities for repentance and seeks to cover and forgive shortcomings, highlighting His compassion and understanding.

**Audit:** Include abundant/repeated forgiveness. Ghaffaar and Ghafoor overlap; do not impose a rigid theological boundary.

- [20:82](https://quran.com/20/82): وَإِنِّي لَغَفَّارٞ لِّمَن تَابَ وَءَامَنَ وَعَمِلَ صَٰلِحٗا ثُمَّ ٱهۡتَدَىٰ
- [38:66](https://quran.com/38/66): رَبُّ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِ وَمَا بَيۡنَهُمَا ٱلۡعَزِيزُ ٱلۡغَفَّٰرُ
- [39:5](https://quran.com/39/5): خَلَقَ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضَ بِٱلۡحَقِّۖ يُكَوِّرُ ٱلَّيۡلَ عَلَى ٱلنَّهَارِ وَيُكَوِّرُ ٱلنَّهَارَ عَلَى ٱلَّيۡلِۖ وَسَخَّرَ ٱلشَّمۡسَ وَٱلۡقَمَرَۖ كُلّٞ يَجۡرِي لِأَجَلٖ مُّسَمًّىۗ أَلَا هُوَ ٱلۡعَزِيزُ ٱلۡغَفَّٰرُ
- [40:42](https://quran.com/40/42): تَدۡعُونَنِي لِأَكۡفُرَ بِٱللَّهِ وَأُشۡرِكَ بِهِۦ مَا لَيۡسَ لِي بِهِۦ عِلۡمٞ وَأَنَا۠ أَدۡعُوكُمۡ إِلَى ٱلۡعَزِيزِ ٱلۡغَفَّٰرِ
- [71:10](https://quran.com/71/10): فَقُلۡتُ ٱسۡتَغۡفِرُواْ رَبَّكُمۡ إِنَّهُۥ كَانَ غَفَّارٗا

### 15. Al Qahhaar — الْقَهَّارُ

**Short label:** The Subduer

**Main description:** Is the One who is almighty to submit all creation to His will.

**Numbered-file description (unused):** Al Qahhaar is one of the names of Allah in Islam, meaning 'The All-Prevailing One' or 'The Supreme Conqueror'. It signifies Allah's ultimate power and control over all creation.

**Displayed summary:** The name Al Qahhaar emphasizes Allah's sovereignty and ability to dominate all beings and situations. It reflects His authority and the inevitable defeat of all that opposes His will. The name inspires believers to submit to Allah's will, reminding them of His might and the transient nature of worldly power.

**Audit:** Subduer and irresistible authority are aligned.

- [13:16](https://quran.com/13/16): قُلۡ مَن رَّبُّ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِ قُلِ ٱللَّهُۚ قُلۡ أَفَٱتَّخَذۡتُم مِّن دُونِهِۦٓ أَوۡلِيَآءَ لَا يَمۡلِكُونَ لِأَنفُسِهِمۡ نَفۡعٗا وَلَا ضَرّٗاۚ قُلۡ هَلۡ يَسۡتَوِي ٱلۡأَعۡمَىٰ وَٱلۡبَصِيرُ أَمۡ هَلۡ تَسۡتَوِي ٱلظُّلُمَٰتُ وَٱلنُّورُۗ أَمۡ جَعَلُواْ لِلَّهِ شُرَكَآءَ خَلَقُواْ كَخَلۡقِهِۦ فَتَشَٰبَهَ ٱلۡخَلۡقُ عَلَيۡهِمۡۚ قُلِ ٱللَّهُ خَٰلِقُ كُلِّ شَيۡءٖ وَهُوَ ٱلۡوَٰحِدُ ٱلۡقَهَّٰرُ
- [14:48](https://quran.com/14/48): يَوۡمَ تُبَدَّلُ ٱلۡأَرۡضُ غَيۡرَ ٱلۡأَرۡضِ وَٱلسَّمَٰوَٰتُۖ وَبَرَزُواْ لِلَّهِ ٱلۡوَٰحِدِ ٱلۡقَهَّارِ
- [38:65](https://quran.com/38/65): قُلۡ إِنَّمَآ أَنَا۠ مُنذِرٞۖ وَمَا مِنۡ إِلَٰهٍ إِلَّا ٱللَّهُ ٱلۡوَٰحِدُ ٱلۡقَهَّارُ
- [39:4](https://quran.com/39/4): لَّوۡ أَرَادَ ٱللَّهُ أَن يَتَّخِذَ وَلَدٗا لَّٱصۡطَفَىٰ مِمَّا يَخۡلُقُ مَا يَشَآءُۚ سُبۡحَٰنَهُۥۖ هُوَ ٱللَّهُ ٱلۡوَٰحِدُ ٱلۡقَهَّارُ
- [40:16](https://quran.com/40/16): يَوۡمَ هُم بَٰرِزُونَۖ لَا يَخۡفَىٰ عَلَى ٱللَّهِ مِنۡهُمۡ شَيۡءٞۚ لِّمَنِ ٱلۡمُلۡكُ ٱلۡيَوۡمَۖ لِلَّهِ ٱلۡوَٰحِدِ ٱلۡقَهَّارِ

### 16. Al Wahhaab — الْوَهَّابُ

**Short label:** The Giver of All

**Main description:** Is the One who gives immensely without receiving anything in return

**Numbered-file description (unused):** Al Wahhaab is one of the names of Allah in Islam, meaning 'The Bestower' or 'The Giver'. It signifies God's generosity and His ability to grant gifts and blessings to His creation without any expectation in return.

**Displayed summary:** The origin of Al Wahhaab comes from the Arabic root word 'wahb', which means to give or to donate freely. This name reflects the aspect of Allah's mercy and benevolence, emphasizing that all provisions and resources come from Him. In Islamic tradition, it is a reminder for believers to be grateful for God's endless gifts.

**Audit:** The Bestower is a concise alternative to Giver of All; retain freely giving bounty.

- [3:8](https://quran.com/3/8): رَبَّنَا لَا تُزِغۡ قُلُوبَنَا بَعۡدَ إِذۡ هَدَيۡتَنَا وَهَبۡ لَنَا مِن لَّدُنكَ رَحۡمَةًۚ إِنَّكَ أَنتَ ٱلۡوَهَّابُ
- [38:9](https://quran.com/38/9): أَمۡ عِندَهُمۡ خَزَآئِنُ رَحۡمَةِ رَبِّكَ ٱلۡعَزِيزِ ٱلۡوَهَّابِ
- [38:35](https://quran.com/38/35): قَالَ رَبِّ ٱغۡفِرۡ لِي وَهَبۡ لِي مُلۡكٗا لَّا يَنۢبَغِي لِأَحَدٖ مِّنۢ بَعۡدِيٓۖ إِنَّكَ أَنتَ ٱلۡوَهَّابُ

### 17. Ar Razzaaq — الرَّزَّاقُ

**Short label:** The Sustainer

**Main description:** Allah provides His creatures with what they need. Provision includes nourishment and other gifts.

**Numbered-file description (unused):** Ar Razzaaq is one of the names of Allah in Islam, meaning 'The Provider' or 'The Sustainer'. It signifies the divine attribute of providing sustenance and means to all of creation.

**Displayed summary:** The name Ar Razzaaq underscores Allah's role as the ultimate provider, ensuring that every creature receives what it needs to sustain life. The origin of this name comes from the Arabic root word 'rizq', which means sustenance. This name reflects the belief that Allah is responsible for the provision of all that is beneficial for His creations, both material and spiritual.

**Audit:** Prefer The Provider, matching provision in 51:58 and the revised description; distinguish wording from Muqeet without denying overlap.

- [51:58](https://quran.com/51/58): إِنَّ ٱللَّهَ هُوَ ٱلرَّزَّاقُ ذُو ٱلۡقُوَّةِ ٱلۡمَتِينُ

### 18. Al Fattaah — الْفَتَّاحُ

**Short label:** The Opener

**Main description:** Is He Who alone knows the keys to the future and the unknown, by His knowledge and by His wisdom. This name is often associated with the following name "Al-‘Alīm".

**Numbered-file description (unused):** Al Fattaah is one of the names of Allah in Islam, meaning 'The Opener' or 'The Victory Giver'. It represents the divine authority to open doors to new opportunities, success, and guidance.

**Displayed summary:** Al Fattaah signifies Allah's role as the One who provides openings and breakthroughs. This name emphasizes the concept of divine intervention in the affairs of humanity, helping believers overcome obstacles and granting them victory. It highlights the importance of seeking Allah's help in times of struggle and the belief that He can change situations for the better.

**Audit:** Summary only describes breakthroughs. Add judging/deciding in truth, the actual context of 34:26.

- [34:26](https://quran.com/34/26): قُلۡ يَجۡمَعُ بَيۡنَنَا رَبُّنَا ثُمَّ يَفۡتَحُ بَيۡنَنَا بِٱلۡحَقِّ وَهُوَ ٱلۡفَتَّاحُ ٱلۡعَلِيمُ

### 19. Al 'Aleem — اَلْعَلِيْمُ

**Short label:** The Knower of all

**Main description:** The one from whom nothing can hide, for He knows everything and is aware of everything.

**Numbered-file description (unused):** Al 'Aleem (العليم) translates to 'The All-Knowing' or 'The Omniscient'. It is one of the names of Allah in Islam, reflecting His complete and perfect knowledge of all things, including the unseen and the seen.

**Displayed summary:** The name Al 'Aleem emphasizes Allah's infinite wisdom and understanding. It signifies that nothing is hidden from His knowledge, and He is aware of every detail in the universe. This name is often invoked in prayers for guidance and understanding, illustrating the importance of knowledge in the Islamic faith.

**Audit:** The All-Knowing is a concise label; knowledge of everything is aligned.

- [2:158](https://quran.com/2/158): ۞ إِنَّ ٱلصَّفَا وَٱلۡمَرۡوَةَ مِن شَعَآئِرِ ٱللَّهِۖ فَمَنۡ حَجَّ ٱلۡبَيۡتَ أَوِ ٱعۡتَمَرَ فَلَا جُنَاحَ عَلَيۡهِ أَن يَطَّوَّفَ بِهِمَاۚ وَمَن تَطَوَّعَ خَيۡرٗا فَإِنَّ ٱللَّهَ شَاكِرٌ عَلِيمٌ
- [3:92](https://quran.com/3/92): لَن تَنَالُواْ ٱلۡبِرَّ حَتَّىٰ تُنفِقُواْ مِمَّا تُحِبُّونَۚ وَمَا تُنفِقُواْ مِن شَيۡءٖ فَإِنَّ ٱللَّهَ بِهِۦ عَلِيمٞ
- [4:35](https://quran.com/4/35): وَإِنۡ خِفۡتُمۡ شِقَاقَ بَيۡنِهِمَا فَٱبۡعَثُواْ حَكَمٗا مِّنۡ أَهۡلِهِۦ وَحَكَمٗا مِّنۡ أَهۡلِهَآ إِن يُرِيدَآ إِصۡلَٰحٗا يُوَفِّقِ ٱللَّهُ بَيۡنَهُمَآۗ إِنَّ ٱللَّهَ كَانَ عَلِيمًا خَبِيرٗا
- [24:41](https://quran.com/24/41): أَلَمۡ تَرَ أَنَّ ٱللَّهَ يُسَبِّحُ لَهُۥ مَن فِي ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِ وَٱلطَّيۡرُ صَٰٓفَّٰتٖۖ كُلّٞ قَدۡ عَلِمَ صَلَاتَهُۥ وَتَسۡبِيحَهُۥۗ وَٱللَّهُ عَلِيمُۢ بِمَا يَفۡعَلُونَ
- [33:40](https://quran.com/33/40): مَّا كَانَ مُحَمَّدٌ أَبَآ أَحَدٖ مِّن رِّجَالِكُمۡ وَلَٰكِن رَّسُولَ ٱللَّهِ وَخَاتَمَ ٱلنَّبِيِّـۧنَۗ وَكَانَ ٱللَّهُ بِكُلِّ شَيۡءٍ عَلِيمٗا

### 20. Al Qaabid — الْقَابِضُ

**Short label:** The Constrictor

**Main description:** Is the One who restricts His good according to wisdom.

**Numbered-file description (unused):** Al Qaabid (Arabic: القابض) is one of the names of Allah in Islamic theology. It translates to 'The Withholder' or 'The Restrainer'.

**Displayed summary:** The name Al Qaabid reflects Allah's attribute of controlling provision and restricting the means of sustenance as He wills. It signifies that Allah has the authority to grant or withhold bounty and that everything is under His control. This name reminds believers of the importance of gratitude and reliance on Allah, as well as the transient nature of worldly possessions.

**Audit:** Quran 2:245 uses a verb, not the listed noun. Pair contraction with expansion; Abu Dawud 3451 explicitly supplies the names.

- [2:245](https://quran.com/2/245): مَّن ذَا ٱلَّذِي يُقۡرِضُ ٱللَّهَ قَرۡضًا حَسَنٗا فَيُضَٰعِفَهُۥ لَهُۥٓ أَضۡعَافٗا كَثِيرَةٗۚ وَٱللَّهُ يَقۡبِضُ وَيَبۡصُۜطُ وَإِلَيۡهِ تُرۡجَعُونَ

### 21. Al Baasit — الْبَاسِطُ

**Short label:** The Reliever

**Main description:** Is the One who frees, increases and multiplies his good according to wisdom.

**Numbered-file description (unused):** Al Baasit is one of the names of Allah, which translates to 'The Expander' or 'The One Who Spreads'. It emphasizes Allah's ability to provide sustenance and abundance, as well as the expansion of mercy and blessings.

**Displayed summary:** The name Al Baasit originates from the Arabic root word 'b-s-t', which conveys the concept of spreading or expanding. In Islamic tradition, it reflects Allah’s role in providing for His creation and ensuring that resources, both spiritual and material, are sufficiently available. This name highlights the boundless nature of Allah's provision and mercy.

**Audit:** Prefer The Expander to The Reliever. 2:245 uses a verb; pair with Qaabid and cite Abu Dawud 3451 separately.

- [2:245](https://quran.com/2/245): مَّن ذَا ٱلَّذِي يُقۡرِضُ ٱللَّهَ قَرۡضًا حَسَنٗا فَيُضَٰعِفَهُۥ لَهُۥٓ أَضۡعَافٗا كَثِيرَةٗۚ وَٱللَّهُ يَقۡبِضُ وَيَبۡصُۜطُ وَإِلَيۡهِ تُرۡجَعُونَ

### 22. Al Khaafid — الْخَافِضُ

**Short label:** The Abaser

**Main description:** Is the One who destines the oppressors and the proud to demeanor and chastisement.

**Numbered-file description (unused):** Al Khaafid (Arabic: الخافض) is one of the names of Allah in Islam, meaning 'The Reducer' or 'The Diminisher'. This name signifies Allah's ability to lower, humiliate, or reduce the status of individuals or things as He wills.

**Displayed summary:** Al Khaafid emphasizes the power of Allah in altering circumstances and statuses, both by raising some and lowering others. The concept reflects His sovereignty and control over all creation, highlighting that one's worth or position is ultimately determined by Him.

**Audit:** Rewrite unclear prose about oppressors and demeanor. 95:5 is a related action, not a literal occurrence of Al Khaafid.

- [95:5](https://quran.com/95/5): ثُمَّ رَدَدۡنَٰهُ أَسۡفَلَ سَٰفِلِينَ

### 23. Ar Raafi' — الرَّافِعُ

**Short label:** The Exalter

**Main description:** The Breeder in the sense of raising, the one who elevates the honor, the status and the prestige of a being in this lower world and in the hereafter.

**Numbered-file description (unused):** Ar Raafi' is one of the names of Allah in Islam that means 'The Exalter' or 'The Elevator'. It signifies Allah's power to raise the status of His creation, be it in their worldly affairs, ranks, or spiritual elevation.

**Displayed summary:** Originating from the Arabic root word 'رفعة' (raf'ah), which means to raise or elevate, Ar Raafi' reflects the belief that Allah can elevate individuals spiritually and socially. This name highlights the attribute of Allah that involves bestowing honor, success, and high status upon those whom He wills.

**Audit:** Remove The Breeder from the description. Both references describe raising through verbs, not the listed name.

- [58:11](https://quran.com/58/11): يَٰٓأَيُّهَا ٱلَّذِينَ ءَامَنُوٓاْ إِذَا قِيلَ لَكُمۡ تَفَسَّحُواْ فِي ٱلۡمَجَٰلِسِ فَٱفۡسَحُواْ يَفۡسَحِ ٱللَّهُ لَكُمۡۖ وَإِذَا قِيلَ ٱنشُزُواْ فَٱنشُزُواْ يَرۡفَعِ ٱللَّهُ ٱلَّذِينَ ءَامَنُواْ مِنكُمۡ وَٱلَّذِينَ أُوتُواْ ٱلۡعِلۡمَ دَرَجَٰتٖۚ وَٱللَّهُ بِمَا تَعۡمَلُونَ خَبِيرٞ
- [6:83](https://quran.com/6/83): وَتِلۡكَ حُجَّتُنَآ ءَاتَيۡنَٰهَآ إِبۡرَٰهِيمَ عَلَىٰ قَوۡمِهِۦۚ نَرۡفَعُ دَرَجَٰتٖ مَّن نَّشَآءُۗ إِنَّ رَبَّكَ حَكِيمٌ عَلِيمٞ

### 24. Al Mu'iz — الْمُعِزُّ

**Short label:** The Bestower of Honour

**Main description:** One who offers honor and precious value.

**Numbered-file description (unused):** Al Mu'iz is one of the names of Allah in Islam, meaning 'The Honorer' or 'The Bestower of Honor.' This name signifies Allah's ability to grant dignity and respect to His creation.

**Displayed summary:** The name Al Mu'iz originates from the Arabic root 'عز' (عزَّ), which conveys themes of honor, strength, and dignity. It reflects Allah's power to elevate believers and provide them with status and respect in the eyes of others and in society. Recognizing Allah as Al Mu'iz encourages humility and reliance on Him as the ultimate source of honor.

**Audit:** Transliteration should reflect doubled z: Al Mu'izz. The Honourer is concise; 3:26 uses a verb.

- [3:26](https://quran.com/3/26): قُلِ ٱللَّهُمَّ مَٰلِكَ ٱلۡمُلۡكِ تُؤۡتِي ٱلۡمُلۡكَ مَن تَشَآءُ وَتَنزِعُ ٱلۡمُلۡكَ مِمَّن تَشَآءُ وَتُعِزُّ مَن تَشَآءُ وَتُذِلُّ مَن تَشَآءُۖ بِيَدِكَ ٱلۡخَيۡرُۖ إِنَّكَ عَلَىٰ كُلِّ شَيۡءٖ قَدِيرٞ

### 25. Al Mudhil — المُذِلُّ

**Short label:** The Humiliator

**Main description:** Is the One who can humiliate and devalue any valiant.

**Numbered-file description (unused):** Al Mudhil is one of the names of Allah in Islam, derived from Arabic. It translates to 'The Dishonorer' or 'The One Who Humiliates.'

**Displayed summary:** Al Mudhil refers to Allah's attribute of causing disgrace or humiliation to those who oppose His commands or reject His guidance. This name emphasizes the power of Allah to elevate or debase individuals based on their actions and faith. It serves as a reminder of the consequences of turning away from righteousness and the importance of adhering to divine guidance.

**Audit:** Transliteration should reflect doubled l: Al Mudhill. Rewrite devalue any valiant; 3:26 uses a verb.

- [3:26](https://quran.com/3/26): قُلِ ٱللَّهُمَّ مَٰلِكَ ٱلۡمُلۡكِ تُؤۡتِي ٱلۡمُلۡكَ مَن تَشَآءُ وَتَنزِعُ ٱلۡمُلۡكَ مِمَّن تَشَآءُ وَتُعِزُّ مَن تَشَآءُ وَتُذِلُّ مَن تَشَآءُۖ بِيَدِكَ ٱلۡخَيۡرُۖ إِنَّكَ عَلَىٰ كُلِّ شَيۡءٖ قَدِيرٞ

### 26. As Samee' — السَّمِيعُ

**Short label:** The Hearer of all

**Main description:** One who hears any noise or sound in His kingdom emanating from any creature and also One who listens to the invocations of His creatures to respond to their invocations.

**Numbered-file description (unused):** As Samee' (Arabic: السميع) is one of the 99 names of Allah in Islam, meaning 'The All-Hearing'. It signifies that Allah hears all things, whether spoken or unspoken, and nothing escapes His hearing.

**Displayed summary:** The name As Samee' emphasizes Allah's attribute of complete awareness and understanding of all sounds and voices. This attribute assures believers that their prayers, thoughts, and cries for help are heard, providing comfort and hope. The origin of this name comes from the Arabic root 'سمع' (sama'), which means 'to hear'. This name reflects Allah's omnipresence and is often invoked in prayers and supplications.

**Audit:** Prefer The All-Hearing. Explain complete hearing without the summary's unnecessary spatial/omnipresence claim.

- [2:127](https://quran.com/2/127): وَإِذۡ يَرۡفَعُ إِبۡرَٰهِـۧمُ ٱلۡقَوَاعِدَ مِنَ ٱلۡبَيۡتِ وَإِسۡمَٰعِيلُ رَبَّنَا تَقَبَّلۡ مِنَّآۖ إِنَّكَ أَنتَ ٱلسَّمِيعُ ٱلۡعَلِيمُ
- [2:256](https://quran.com/2/256): لَآ إِكۡرَاهَ فِي ٱلدِّينِۖ قَد تَّبَيَّنَ ٱلرُّشۡدُ مِنَ ٱلۡغَيِّۚ فَمَن يَكۡفُرۡ بِٱلطَّٰغُوتِ وَيُؤۡمِنۢ بِٱللَّهِ فَقَدِ ٱسۡتَمۡسَكَ بِٱلۡعُرۡوَةِ ٱلۡوُثۡقَىٰ لَا ٱنفِصَامَ لَهَاۗ وَٱللَّهُ سَمِيعٌ عَلِيمٌ
- [8:17](https://quran.com/8/17): فَلَمۡ تَقۡتُلُوهُمۡ وَلَٰكِنَّ ٱللَّهَ قَتَلَهُمۡۚ وَمَا رَمَيۡتَ إِذۡ رَمَيۡتَ وَلَٰكِنَّ ٱللَّهَ رَمَىٰ وَلِيُبۡلِيَ ٱلۡمُؤۡمِنِينَ مِنۡهُ بَلَآءً حَسَنًاۚ إِنَّ ٱللَّهَ سَمِيعٌ عَلِيمٞ
- [49:1](https://quran.com/49/1): يَٰٓأَيُّهَا ٱلَّذِينَ ءَامَنُواْ لَا تُقَدِّمُواْ بَيۡنَ يَدَيِ ٱللَّهِ وَرَسُولِهِۦۖ وَٱتَّقُواْ ٱللَّهَۚ إِنَّ ٱللَّهَ سَمِيعٌ عَلِيمٞ

### 27. Al Baseer — الْبَصِيرُ

**Short label:** The Seer of all

**Main description:** Is the One from whom nothing that exists or comes into existence can escape without the slightest use of an eye, instrument or organ.

**Numbered-file description (unused):** Al Baseer, meaning 'The All-Seeing', is one of the names of Allah in Islam that signifies His infinite ability to perceive and understand everything in existence.

**Displayed summary:** The name Al Baseer highlights Allah's attribute of seeing all things, whether visible or hidden. It represents His knowledge of the inner thoughts and intentions of individuals, as well as the external realities of the world. Originating from the Arabic root 'بصر', which means sight or vision, Al Baseer reassures believers that nothing escapes Allah's awareness, encouraging them to act justly and with sincerity.

**Audit:** Prefer The All-Seeing; seeing all things is aligned.

- [4:58](https://quran.com/4/58): ۞ إِنَّ ٱللَّهَ يَأۡمُرُكُمۡ أَن تُؤَدُّواْ ٱلۡأَمَٰنَٰتِ إِلَىٰٓ أَهۡلِهَا وَإِذَا حَكَمۡتُم بَيۡنَ ٱلنَّاسِ أَن تَحۡكُمُواْ بِٱلۡعَدۡلِۚ إِنَّ ٱللَّهَ نِعِمَّا يَعِظُكُم بِهِۦٓۗ إِنَّ ٱللَّهَ كَانَ سَمِيعَۢا بَصِيرٗا
- [17:1](https://quran.com/17/1): سُبۡحَٰنَ ٱلَّذِيٓ أَسۡرَىٰ بِعَبۡدِهِۦ لَيۡلٗا مِّنَ ٱلۡمَسۡجِدِ ٱلۡحَرَامِ إِلَى ٱلۡمَسۡجِدِ ٱلۡأَقۡصَا ٱلَّذِي بَٰرَكۡنَا حَوۡلَهُۥ لِنُرِيَهُۥ مِنۡ ءَايَٰتِنَآۚ إِنَّهُۥ هُوَ ٱلسَّمِيعُ ٱلۡبَصِيرُ
- [42:11](https://quran.com/42/11): فَاطِرُ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِۚ جَعَلَ لَكُم مِّنۡ أَنفُسِكُمۡ أَزۡوَٰجٗا وَمِنَ ٱلۡأَنۡعَٰمِ أَزۡوَٰجٗا يَذۡرَؤُكُمۡ فِيهِۚ لَيۡسَ كَمِثۡلِهِۦ شَيۡءٞۖ وَهُوَ ٱلسَّمِيعُ ٱلۡبَصِيرُ
- [42:27](https://quran.com/42/27): ۞ وَلَوۡ بَسَطَ ٱللَّهُ ٱلرِّزۡقَ لِعِبَادِهِۦ لَبَغَوۡاْ فِي ٱلۡأَرۡضِ وَلَٰكِن يُنَزِّلُ بِقَدَرٖ مَّا يَشَآءُۚ إِنَّهُۥ بِعِبَادِهِۦ خَبِيرُۢ بَصِيرٞ

### 28. Al Hakam — الْحَكَمُ

**Short label:** The Judge

**Main description:** He is the One who judges and decides between creatures here below and in the hereafter and He who ensures equity between the creatures and no one else will and will not be able to decide between them and judge them; and the Judgment in the hereafter will be perfectly just and appropriate.

**Numbered-file description (unused):** Al Hakam is one of the names of Allah in Islamic tradition, meaning 'The Judge' or 'The Arbitrator'.

**Displayed summary:** The name Al Hakam emphasizes Allah's role as the supreme judge who determines justice and fairness. It signifies His authority to rule over His creation and to decide matters, rewarding the righteous and punishing the wrongdoers. The use of this name serves to remind believers of the ultimate accountability they have in the Hereafter.

**Audit:** Shorten repetitive description; add 6:114, which uses hakaman for Allah. Existing 22:69 uses the verb judges.

- [22:69](https://quran.com/22/69): ٱللَّهُ يَحۡكُمُ بَيۡنَكُمۡ يَوۡمَ ٱلۡقِيَٰمَةِ فِيمَا كُنتُمۡ فِيهِ تَخۡتَلِفُونَ

### 29. Al 'Adl — الْعَدْلُ

**Short label:** The Just

**Main description:** Is the One who is free from all forms of injustice and / or oppression.

**Numbered-file description (unused):** Al 'Adl is one of the 99 Names of Allah in Islam, meaning 'The Just' or 'The Utterly Just'. It signifies God's inherent justice and fairness in all matters.

**Displayed summary:** The concept of Al 'Adl reflects the belief that Allah is just in His actions, laws, and decrees, ensuring fairness in the universe. This name emphasizes the importance of justice as a moral value in Islam, guiding individual behavior and societal laws. The origin of Al 'Adl is rooted in Islamic theology, where it underscores the belief that Allah does not oppress anyone and that all human beings will be judged fairly on the Day of Judgment.

**Audit:** 6:115 describes the justice of Allah's word, not a literal Al Adl. 4:40 is a clearer related passage for absence of injustice.

- [6:115](https://quran.com/6/115): وَتَمَّتۡ كَلِمَتُ رَبِّكَ صِدۡقٗا وَعَدۡلٗاۚ لَّا مُبَدِّلَ لِكَلِمَٰتِهِۦۚ وَهُوَ ٱلسَّمِيعُ ٱلۡعَلِيمُ

### 30. Al Lateef — اللَّطِيفُ

**Short label:** The Subtle One

**Main description:** Is the One who grants his creatures His gentleness and His benevolence.

**Numbered-file description (unused):** Al Lateef (اللَّطِيفُ) is one of the names of Allah in Islam, meaning 'The Subtle One' or 'The Gracious'. It signifies Allah's gentle and kind nature, as well as His ability to provide for His creation in ways that may not be immediately apparent.

**Displayed summary:** The name Al Lateef originates from the Arabic root 'l-t-f', which conveys the concept of kindness and subtlety. In Islamic theology, it emphasizes Allah's care and gentleness, reminding believers of His hidden grace and support in their lives. Al Lateef assures followers that even in difficult times, there is a divine presence that understands and provides for their needs with kindness.

**Audit:** Explain both subtle knowledge and gracious kindness; the current description emphasises only gentleness.

- [6:103](https://quran.com/6/103): لَّا تُدۡرِكُهُ ٱلۡأَبۡصَٰرُ وَهُوَ يُدۡرِكُ ٱلۡأَبۡصَٰرَۖ وَهُوَ ٱللَّطِيفُ ٱلۡخَبِيرُ
- [22:63](https://quran.com/22/63): أَلَمۡ تَرَ أَنَّ ٱللَّهَ أَنزَلَ مِنَ ٱلسَّمَآءِ مَآءٗ فَتُصۡبِحُ ٱلۡأَرۡضُ مُخۡضَرَّةًۚ إِنَّ ٱللَّهَ لَطِيفٌ خَبِيرٞ
- [31:16](https://quran.com/31/16): يَٰبُنَيَّ إِنَّهَآ إِن تَكُ مِثۡقَالَ حَبَّةٖ مِّنۡ خَرۡدَلٖ فَتَكُن فِي صَخۡرَةٍ أَوۡ فِي ٱلسَّمَٰوَٰتِ أَوۡ فِي ٱلۡأَرۡضِ يَأۡتِ بِهَا ٱللَّهُۚ إِنَّ ٱللَّهَ لَطِيفٌ خَبِيرٞ
- [33:34](https://quran.com/33/34): وَٱذۡكُرۡنَ مَا يُتۡلَىٰ فِي بُيُوتِكُنَّ مِنۡ ءَايَٰتِ ٱللَّهِ وَٱلۡحِكۡمَةِۚ إِنَّ ٱللَّهَ كَانَ لَطِيفًا خَبِيرًا

### 31. Al Khabeer — الْخَبِيرُ

**Short label:** The All Aware

**Main description:** He is the One from whom cannot be hidden the reality of everything that exists, that is to say the truth of things; no detail escapes Him, neither the thing as a whole, nor each of its parts; and to claim to be able to hide things from the One who creates everything is a delusion and a disbelief.

**Numbered-file description (unused):** Al Khabeer, which translates to 'The All-Aware' or 'The All-Knowing', is one of the names of Allah in Islam. It signifies God's complete and profound knowledge of all things, including the hidden and the apparent. This attribute emphasizes His awareness of the intricacies of creation and human actions.

**Displayed summary:** The name Al Khabeer comes from the Arabic root 'k-b-r', related to deep knowledge and awareness. It is a reminder that Allah possesses knowledge beyond human comprehension, encompassing the unseen and the seen. The name encourages believers to seek Allah's guidance and trust in His wisdom.

**Audit:** Correct summary root k-b-r to kh-b-r (خ ب ر). Remove unrelated moralising from the main description.

- [6:18](https://quran.com/6/18): وَهُوَ ٱلۡقَاهِرُ فَوۡقَ عِبَادِهِۦۚ وَهُوَ ٱلۡحَكِيمُ ٱلۡخَبِيرُ
- [17:30](https://quran.com/17/30): إِنَّ رَبَّكَ يَبۡسُطُ ٱلرِّزۡقَ لِمَن يَشَآءُ وَيَقۡدِرُۚ إِنَّهُۥ كَانَ بِعِبَادِهِۦ خَبِيرَۢا بَصِيرٗا
- [49:13](https://quran.com/49/13): يَٰٓأَيُّهَا ٱلنَّاسُ إِنَّا خَلَقۡنَٰكُم مِّن ذَكَرٖ وَأُنثَىٰ وَجَعَلۡنَٰكُمۡ شُعُوبٗا وَقَبَآئِلَ لِتَعَارَفُوٓاْۚ إِنَّ أَكۡرَمَكُمۡ عِندَ ٱللَّهِ أَتۡقَىٰكُمۡۚ إِنَّ ٱللَّهَ عَلِيمٌ خَبِيرٞ
- [59:18](https://quran.com/59/18): يَٰٓأَيُّهَا ٱلَّذِينَ ءَامَنُواْ ٱتَّقُواْ ٱللَّهَ وَلۡتَنظُرۡ نَفۡسٞ مَّا قَدَّمَتۡ لِغَدٖۖ وَٱتَّقُواْ ٱللَّهَۚ إِنَّ ٱللَّهَ خَبِيرُۢ بِمَا تَعۡمَلُونَ

### 32. Al Haleem — الْحَلِيمُ

**Short label:** The Forbearing

**Main description:** Is the One whose patience towards the disbeliefs of sinners exceeds all imagination (it is possible that forgiveness will be granted to them).

**Numbered-file description (unused):** Al Haleem is one of the names of Allah in Islam, meaning 'The Forbearing' or 'The Patient'. It signifies the quality of being slow to punish and tolerant in the face of the shortcomings of others.

**Displayed summary:** The name Al Haleem originates from the Arabic root 'H-L-M', which conveys the meaning of gentleness and forbearance. In the Quran, Allah is described as Al Haleem, emphasizing His patience and restraint towards His creation, even when they commit sins. This attribute encourages believers to embody patience and understanding in their own lives.

**Audit:** Keep The Forbearing; simplify to not hastening punishment despite knowing wrongdoing.

- [2:235](https://quran.com/2/235): وَلَا جُنَاحَ عَلَيۡكُمۡ فِيمَا عَرَّضۡتُم بِهِۦ مِنۡ خِطۡبَةِ ٱلنِّسَآءِ أَوۡ أَكۡنَنتُمۡ فِيٓ أَنفُسِكُمۡۚ عَلِمَ ٱللَّهُ أَنَّكُمۡ سَتَذۡكُرُونَهُنَّ وَلَٰكِن لَّا تُوَاعِدُوهُنَّ سِرًّا إِلَّآ أَن تَقُولُواْ قَوۡلٗا مَّعۡرُوفٗاۚ وَلَا تَعۡزِمُواْ عُقۡدَةَ ٱلنِّكَاحِ حَتَّىٰ يَبۡلُغَ ٱلۡكِتَٰبُ أَجَلَهُۥۚ وَٱعۡلَمُوٓاْ أَنَّ ٱللَّهَ يَعۡلَمُ مَا فِيٓ أَنفُسِكُمۡ فَٱحۡذَرُوهُۚ وَٱعۡلَمُوٓاْ أَنَّ ٱللَّهَ غَفُورٌ حَلِيمٞ
- [17:44](https://quran.com/17/44): تُسَبِّحُ لَهُ ٱلسَّمَٰوَٰتُ ٱلسَّبۡعُ وَٱلۡأَرۡضُ وَمَن فِيهِنَّۚ وَإِن مِّن شَيۡءٍ إِلَّا يُسَبِّحُ بِحَمۡدِهِۦ وَلَٰكِن لَّا تَفۡقَهُونَ تَسۡبِيحَهُمۡۚ إِنَّهُۥ كَانَ حَلِيمًا غَفُورٗا
- [22:59](https://quran.com/22/59): لَيُدۡخِلَنَّهُم مُّدۡخَلٗا يَرۡضَوۡنَهُۥۚ وَإِنَّ ٱللَّهَ لَعَلِيمٌ حَلِيمٞ
- [35:41](https://quran.com/35/41): ۞ إِنَّ ٱللَّهَ يُمۡسِكُ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضَ أَن تَزُولَاۚ وَلَئِن زَالَتَآ إِنۡ أَمۡسَكَهُمَا مِنۡ أَحَدٖ مِّنۢ بَعۡدِهِۦٓۚ إِنَّهُۥ كَانَ حَلِيمًا غَفُورٗا

### 33. Al 'Azeem — الْعَظِيمُ

**Short label:** The Magnificent

**Main description:** This name signifies the absolute and immeasurable greatness of God. It implies a level of grandeur and majesty that is beyond human comprehension.

**Numbered-file description (unused):** Al 'Azeem is one of the names of Allah in Islam, meaning 'The Most Great' or 'The Supreme.'

**Displayed summary:** This name reflects Allah's greatness and majesty. It signifies His unparalleled power and authority over all creation. The name Al 'Azeem emphasizes that no one and nothing can compare to Allah in greatness. It originates from the Arabic root word 'عَظَمَ' (ʿaẓama), which means to be great or to exalt. In prayers and supplications, invoking this name serves as a reminder of Allah's immense greatness and serves as an expression of humility and awe from the believer.

**Audit:** Magnificence is aligned; avoid manufacturing a strict separation from Kabeer.

- [2:255](https://quran.com/2/255): ٱللَّهُ لَآ إِلَٰهَ إِلَّا هُوَ ٱلۡحَيُّ ٱلۡقَيُّومُۚ لَا تَأۡخُذُهُۥ سِنَةٞ وَلَا نَوۡمٞۚ لَّهُۥ مَا فِي ٱلسَّمَٰوَٰتِ وَمَا فِي ٱلۡأَرۡضِۗ مَن ذَا ٱلَّذِي يَشۡفَعُ عِندَهُۥٓ إِلَّا بِإِذۡنِهِۦۚ يَعۡلَمُ مَا بَيۡنَ أَيۡدِيهِمۡ وَمَا خَلۡفَهُمۡۖ وَلَا يُحِيطُونَ بِشَيۡءٖ مِّنۡ عِلۡمِهِۦٓ إِلَّا بِمَا شَآءَۚ وَسِعَ كُرۡسِيُّهُ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضَۖ وَلَا يَـُٔودُهُۥ حِفۡظُهُمَاۚ وَهُوَ ٱلۡعَلِيُّ ٱلۡعَظِيمُ
- [42:4](https://quran.com/42/4): لَهُۥ مَا فِي ٱلسَّمَٰوَٰتِ وَمَا فِي ٱلۡأَرۡضِۖ وَهُوَ ٱلۡعَلِيُّ ٱلۡعَظِيمُ
- [56:96](https://quran.com/56/96): فَسَبِّحۡ بِٱسۡمِ رَبِّكَ ٱلۡعَظِيمِ

### 34. Al Ghafoor — الْغَفُورُ

**Short label:** The Great Forgiver

**Main description:** Is the One who forgives immensely over and over again.

**Numbered-file description (unused):** Al Ghafoor is one of the names of Allah in Islam, meaning 'The Most Forgiving'. It signifies Allah's capacity to forgive sins and show mercy to His creation.

**Displayed summary:** The name Al Ghafoor originates from the Arabic root 'gh-f-r', which means to cover or to conceal. In the context of Allah's attributes, it reflects His ability to overlook and forgive the mistakes and wrongdoings of His servants, encouraging repentance and hope for mercy. Al Ghafoor emphasizes Allah's compassion and readiness to pardon, reinforcing the importance of seeking forgiveness in one's relationship with Him.

**Audit:** Forgiveness is aligned; make its overlap with Ghaffaar explicit rather than inventing separate categories of sins.

- [2:173](https://quran.com/2/173): إِنَّمَا حَرَّمَ عَلَيۡكُمُ ٱلۡمَيۡتَةَ وَٱلدَّمَ وَلَحۡمَ ٱلۡخِنزِيرِ وَمَآ أُهِلَّ بِهِۦ لِغَيۡرِ ٱللَّهِۖ فَمَنِ ٱضۡطُرَّ غَيۡرَ بَاغٖ وَلَا عَادٖ فَلَآ إِثۡمَ عَلَيۡهِۚ إِنَّ ٱللَّهَ غَفُورٞ رَّحِيمٌ
- [8:69](https://quran.com/8/69): فَكُلُواْ مِمَّا غَنِمۡتُمۡ حَلَٰلٗا طَيِّبٗاۚ وَٱتَّقُواْ ٱللَّهَۚ إِنَّ ٱللَّهَ غَفُورٞ رَّحِيمٞ
- [16:110](https://quran.com/16/110): ثُمَّ إِنَّ رَبَّكَ لِلَّذِينَ هَاجَرُواْ مِنۢ بَعۡدِ مَا فُتِنُواْ ثُمَّ جَٰهَدُواْ وَصَبَرُوٓاْ إِنَّ رَبَّكَ مِنۢ بَعۡدِهَا لَغَفُورٞ رَّحِيمٞ
- [41:32](https://quran.com/41/32): نُزُلٗا مِّنۡ غَفُورٖ رَّحِيمٖ

### 35. Ash Shakoor — الشَّكُورُ

**Short label:** The Rewarder of Thankfulness

**Main description:** Is the One who shows His gratitude for worshipers, however incredibly easy to perform.

**Numbered-file description (unused):** Ash Shakoor is one of the names of Allah in Islam, meaning 'The Appreciative' or 'The Acknowledging'. It signifies Allah's recognition and reward for the good deeds of His servants.

**Displayed summary:** The term 'Ash Shakoor' originates from Arabic, where 'Shakoor' refers to being grateful or appreciative. This name reflects Allah's nature of valuing and rewarding the efforts and sincerity of His worshippers, even if their deeds may seem minor. It emphasizes the importance of gratitude and acknowledgment, encouraging believers to perform good deeds with the understanding that Allah appreciates their efforts.

**Audit:** Prefer The Appreciative; explain generous reward for sincere good rather than only rewarding thankfulness. Repair awkward description.

- [35:30](https://quran.com/35/30): لِيُوَفِّيَهُمۡ أُجُورَهُمۡ وَيَزِيدَهُم مِّن فَضۡلِهِۦٓۚ إِنَّهُۥ غَفُورٞ شَكُورٞ
- [35:34](https://quran.com/35/34): وَقَالُواْ ٱلۡحَمۡدُ لِلَّهِ ٱلَّذِيٓ أَذۡهَبَ عَنَّا ٱلۡحَزَنَۖ إِنَّ رَبَّنَا لَغَفُورٞ شَكُورٌ
- [42:23](https://quran.com/42/23): ذَٰلِكَ ٱلَّذِي يُبَشِّرُ ٱللَّهُ عِبَادَهُ ٱلَّذِينَ ءَامَنُواْ وَعَمِلُواْ ٱلصَّٰلِحَٰتِۗ قُل لَّآ أَسۡـَٔلُكُمۡ عَلَيۡهِ أَجۡرًا إِلَّا ٱلۡمَوَدَّةَ فِي ٱلۡقُرۡبَىٰۗ وَمَن يَقۡتَرِفۡ حَسَنَةٗ نَّزِدۡ لَهُۥ فِيهَا حُسۡنًاۚ إِنَّ ٱللَّهَ غَفُورٞ شَكُورٌ
- [64:17](https://quran.com/64/17): إِن تُقۡرِضُواْ ٱللَّهَ قَرۡضًا حَسَنٗا يُضَٰعِفۡهُ لَكُمۡ وَيَغۡفِرۡ لَكُمۡۚ وَٱللَّهُ شَكُورٌ حَلِيمٌ

### 36. Al 'Aliyy — الْعَلِيُّ

**Short label:** The Highest

**Main description:** Is the One who is above all that He created perfectly and unimaginably.

**Numbered-file description (unused):** Al 'Aliyy is one of the names of Allah in Islam, meaning 'The Most High' or 'The Exalted'.

**Displayed summary:** The name Al 'Aliyy emphasizes the supreme status of Allah, highlighting His greatness, majesty, and unparalleled position above all creation. It reflects His authority and power over the universe, as well as His ability to elevate and honor whomever He wills. This name is derived from the Arabic root 'ع ل ي', which connotes highness and elevation.

**Audit:** Exaltedness is aligned; keep explanation clear without adding a disputed spatial interpretation.

- [4:34](https://quran.com/4/34): ٱلرِّجَالُ قَوَّٰمُونَ عَلَى ٱلنِّسَآءِ بِمَا فَضَّلَ ٱللَّهُ بَعۡضَهُمۡ عَلَىٰ بَعۡضٖ وَبِمَآ أَنفَقُواْ مِنۡ أَمۡوَٰلِهِمۡۚ فَٱلصَّٰلِحَٰتُ قَٰنِتَٰتٌ حَٰفِظَٰتٞ لِّلۡغَيۡبِ بِمَا حَفِظَ ٱللَّهُۚ وَٱلَّٰتِي تَخَافُونَ نُشُوزَهُنَّ فَعِظُوهُنَّ وَٱهۡجُرُوهُنَّ فِي ٱلۡمَضَاجِعِ وَٱضۡرِبُوهُنَّۖ فَإِنۡ أَطَعۡنَكُمۡ فَلَا تَبۡغُواْ عَلَيۡهِنَّ سَبِيلًاۗ إِنَّ ٱللَّهَ كَانَ عَلِيّٗا كَبِيرٗا
- [31:30](https://quran.com/31/30): ذَٰلِكَ بِأَنَّ ٱللَّهَ هُوَ ٱلۡحَقُّ وَأَنَّ مَا يَدۡعُونَ مِن دُونِهِ ٱلۡبَٰطِلُ وَأَنَّ ٱللَّهَ هُوَ ٱلۡعَلِيُّ ٱلۡكَبِيرُ
- [42:4](https://quran.com/42/4): لَهُۥ مَا فِي ٱلسَّمَٰوَٰتِ وَمَا فِي ٱلۡأَرۡضِۖ وَهُوَ ٱلۡعَلِيُّ ٱلۡعَظِيمُ
- [42:51](https://quran.com/42/51): ۞ وَمَا كَانَ لِبَشَرٍ أَن يُكَلِّمَهُ ٱللَّهُ إِلَّا وَحۡيًا أَوۡ مِن وَرَآيِٕ حِجَابٍ أَوۡ يُرۡسِلَ رَسُولٗا فَيُوحِيَ بِإِذۡنِهِۦ مَا يَشَآءُۚ إِنَّهُۥ عَلِيٌّ حَكِيمٞ

### 37. Al Kabeer — الْكَبِيرُ

**Short label:** The Greatest

**Main description:** His greatness is the perfection of His being: His existence is eternal, and all created existence depends on Him.

**Numbered-file description (unused):** Al Kabeer is one of the names of Allah in Islam, which translates to 'The Most Great' or 'The Most Exalted.' It emphasizes God's greatness and supremacy over all creation.

**Displayed summary:** The term Al Kabeer highlights Allah's unparalleled magnitude, authority, and grandeur, reminding believers of His infinite power and majesty. The name originates from the Arabic root 'K-B-R,' meaning 'to be great or to be large.' It is found in various Islamic texts, emphasizing the significance of recognizing Allah's greatness in worship and daily life.

**Audit:** The Greatest is aligned with the revised description; meaningful overlap with other greatness names is expected.

- [13:9](https://quran.com/13/9): عَٰلِمُ ٱلۡغَيۡبِ وَٱلشَّهَٰدَةِ ٱلۡكَبِيرُ ٱلۡمُتَعَالِ
- [22:62](https://quran.com/22/62): ذَٰلِكَ بِأَنَّ ٱللَّهَ هُوَ ٱلۡحَقُّ وَأَنَّ مَا يَدۡعُونَ مِن دُونِهِۦ هُوَ ٱلۡبَٰطِلُ وَأَنَّ ٱللَّهَ هُوَ ٱلۡعَلِيُّ ٱلۡكَبِيرُ
- [31:30](https://quran.com/31/30): ذَٰلِكَ بِأَنَّ ٱللَّهَ هُوَ ٱلۡحَقُّ وَأَنَّ مَا يَدۡعُونَ مِن دُونِهِ ٱلۡبَٰطِلُ وَأَنَّ ٱللَّهَ هُوَ ٱلۡعَلِيُّ ٱلۡكَبِيرُ

### 38. Al Hafeez — الْحَفِيظُ

**Short label:** The Preserver

**Main description:** Is the One who protects and preserves His good.

**Numbered-file description (unused):** Al Hafeez is one of the names of Allah in Islam, meaning 'The Preserver' or 'The Guardian'.

**Displayed summary:** The name Al Hafeez emphasizes Allah's attribute of protecting and preserving all of creation, ensuring the safety of the universe and its elements. It derives from the Arabic root 'H-f-z', which implies safeguarding and guardianship. This name reflects the belief that Allah is vigilant over His creation, preserving their essence, fate, and the balance of the cosmos.

**Audit:** Clarify the vague preserves His good wording: preserving/guarding and knowledge of deeds.

- [11:57](https://quran.com/11/57): فَإِن تَوَلَّوۡاْ فَقَدۡ أَبۡلَغۡتُكُم مَّآ أُرۡسِلۡتُ بِهِۦٓ إِلَيۡكُمۡۚ وَيَسۡتَخۡلِفُ رَبِّي قَوۡمًا غَيۡرَكُمۡ وَلَا تَضُرُّونَهُۥ شَيۡـًٔاۚ إِنَّ رَبِّي عَلَىٰ كُلِّ شَيۡءٍ حَفِيظٞ
- [34:21](https://quran.com/34/21): وَمَا كَانَ لَهُۥ عَلَيۡهِم مِّن سُلۡطَٰنٍ إِلَّا لِنَعۡلَمَ مَن يُؤۡمِنُ بِٱلۡأٓخِرَةِ مِمَّنۡ هُوَ مِنۡهَا فِي شَكّٖۗ وَرَبُّكَ عَلَىٰ كُلِّ شَيۡءٍ حَفِيظٞ
- [42:6](https://quran.com/42/6): وَٱلَّذِينَ ٱتَّخَذُواْ مِن دُونِهِۦٓ أَوۡلِيَآءَ ٱللَّهُ حَفِيظٌ عَلَيۡهِمۡ وَمَآ أَنتَ عَلَيۡهِم بِوَكِيلٖ

### 39. Al Muqeet — المُقيِت

**Short label:** The Sustainer

**Main description:** Allah supplies nourishment to bodies and hearts. In this sense, Al Muqeet is more specific than Ar Razzaaq, whose provision includes nourishment and other gifts.

**Numbered-file description (unused):** Al Muqeet is one of the names of Allah in Islam, meaning 'The Sustainer' or 'The Provider'. It signifies Allah's role in providing and maintaining the necessities of life for His creation.

**Displayed summary:** The name Al Muqeet comes from the Arabic root 'q-t', which relates to provision and sustenance. It reflects Allah's capacity to care for and nourish all beings, ensuring they receive what they need to survive and thrive, both physically and spiritually. This name emphasizes divine mercy and support.

**Audit:** The Sustainer fits the existing revised description. Correct incomplete q-t root to q-w-t (ق و ت); acknowledge interpretations involving power/watchfulness in 4:85.

- [4:85](https://quran.com/4/85): مَّن يَشۡفَعۡ شَفَٰعَةً حَسَنَةٗ يَكُن لَّهُۥ نَصِيبٞ مِّنۡهَاۖ وَمَن يَشۡفَعۡ شَفَٰعَةٗ سَيِّئَةٗ يَكُن لَّهُۥ كِفۡلٞ مِّنۡهَاۗ وَكَانَ ٱللَّهُ عَلَىٰ كُلِّ شَيۡءٖ مُّقِيتٗا

### 40. Al Haseeb — الْحسِيبُ

**Short label:** The Reckoner, The Sufficient

**Main description:** He is the One who will count all the facts of His creatures without forgetting or error.

**Numbered-file description (unused):** Al Haseeb is one of the names of Allah in Islam, meaning 'The Reckoner' or 'The One Who Keeps Account'.

**Displayed summary:** The name Al Haseeb signifies that Allah is the All-Sufficient One who takes account of all actions. It reflects His attribute of justice and emphasizes that He is fully aware of everything that occurs in the universe. This name conveys the idea that individuals will be held accountable for their deeds, encouraging personal responsibility and righteousness in behavior.

**Audit:** Use The Reckoner as the short label and retain sufficiency in the fuller explanation; avoid a multi-label heading.

- [4:6](https://quran.com/4/6): وَٱبۡتَلُواْ ٱلۡيَتَٰمَىٰ حَتَّىٰٓ إِذَا بَلَغُواْ ٱلنِّكَاحَ فَإِنۡ ءَانَسۡتُم مِّنۡهُمۡ رُشۡدٗا فَٱدۡفَعُوٓاْ إِلَيۡهِمۡ أَمۡوَٰلَهُمۡۖ وَلَا تَأۡكُلُوهَآ إِسۡرَافٗا وَبِدَارًا أَن يَكۡبَرُواْۚ وَمَن كَانَ غَنِيّٗا فَلۡيَسۡتَعۡفِفۡۖ وَمَن كَانَ فَقِيرٗا فَلۡيَأۡكُلۡ بِٱلۡمَعۡرُوفِۚ فَإِذَا دَفَعۡتُمۡ إِلَيۡهِمۡ أَمۡوَٰلَهُمۡ فَأَشۡهِدُواْ عَلَيۡهِمۡۚ وَكَفَىٰ بِٱللَّهِ حَسِيبٗا
- [4:86](https://quran.com/4/86): وَإِذَا حُيِّيتُم بِتَحِيَّةٖ فَحَيُّواْ بِأَحۡسَنَ مِنۡهَآ أَوۡ رُدُّوهَآۗ إِنَّ ٱللَّهَ كَانَ عَلَىٰ كُلِّ شَيۡءٍ حَسِيبًا
- [33:39](https://quran.com/33/39): ٱلَّذِينَ يُبَلِّغُونَ رِسَٰلَٰتِ ٱللَّهِ وَيَخۡشَوۡنَهُۥ وَلَا يَخۡشَوۡنَ أَحَدًا إِلَّا ٱللَّهَۗ وَكَفَىٰ بِٱللَّهِ حَسِيبٗا

### 41. Al Jaleel — الْجَلِيلُ

**Short label:** The Majestic

**Main description:** Is the One who is perfectly exempt from being compared to His creatures or from being like Him.

**Numbered-file description (unused):** Al Jaleel is one of the names of Allah in Islam, meaning 'The Majestic' or 'The Glorious'. It signifies the immense greatness and dignity of Allah.

**Displayed summary:** The name Al Jaleel emphasizes Allah's majesty and grandeur, reflecting His supreme authority and the awe-inspiring nature of His creation. It reminds believers of the reverence due to Allah and the beauty in His divine attributes. The name is derived from the Arabic root 'jalal', which relates to honor and magnificence.

**Audit:** 55:27 contains Dhu al-Jalal, a related expression, not Al Jaleel. Remove 39:14 as unrelated; 7:143 is thematic, not a name occurrence.

- [55:27](https://quran.com/55/27): وَيَبۡقَىٰ وَجۡهُ رَبِّكَ ذُو ٱلۡجَلَٰلِ وَٱلۡإِكۡرَامِ
- [39:14](https://quran.com/39/14): قُلِ ٱللَّهَ أَعۡبُدُ مُخۡلِصٗا لَّهُۥ دِينِي
- [7:143](https://quran.com/7/143): وَلَمَّا جَآءَ مُوسَىٰ لِمِيقَٰتِنَا وَكَلَّمَهُۥ رَبُّهُۥ قَالَ رَبِّ أَرِنِيٓ أَنظُرۡ إِلَيۡكَۚ قَالَ لَن تَرَىٰنِي وَلَٰكِنِ ٱنظُرۡ إِلَى ٱلۡجَبَلِ فَإِنِ ٱسۡتَقَرَّ مَكَانَهُۥ فَسَوۡفَ تَرَىٰنِيۚ فَلَمَّا تَجَلَّىٰ رَبُّهُۥ لِلۡجَبَلِ جَعَلَهُۥ دَكّٗا وَخَرَّ مُوسَىٰ صَعِقٗاۚ فَلَمَّآ أَفَاقَ قَالَ سُبۡحَٰنَكَ تُبۡتُ إِلَيۡكَ وَأَنَا۠ أَوَّلُ ٱلۡمُؤۡمِنِينَ

### 42. Al Kareem — الْكَرِيمُ

**Short label:** The Most Generous

**Main description:** Is the One who flooded His creation with His generosity.

**Numbered-file description (unused):** Al-Kareem is one of the names of Allah in Islam, which translates to 'The Generous' or 'The Bountiful'. It signifies God's abundant grace and generosity towards His creation.

**Displayed summary:** Originating from Arabic, Al-Kareem is derived from the root word 'k-r-m' which means to be generous or noble. In Islamic theology, this name emphasizes Allah's boundless mercy and His willingness to bestow blessings and sustenance upon His creatures without any expectation in return.

**Audit:** Generosity is aligned; keep short label and explanation consistent.

- [27:40](https://quran.com/27/40): قَالَ ٱلَّذِي عِندَهُۥ عِلۡمٞ مِّنَ ٱلۡكِتَٰبِ أَنَا۠ ءَاتِيكَ بِهِۦ قَبۡلَ أَن يَرۡتَدَّ إِلَيۡكَ طَرۡفُكَۚ فَلَمَّا رَءَاهُ مُسۡتَقِرًّا عِندَهُۥ قَالَ هَٰذَا مِن فَضۡلِ رَبِّي لِيَبۡلُوَنِيٓ ءَأَشۡكُرُ أَمۡ أَكۡفُرُۖ وَمَن شَكَرَ فَإِنَّمَا يَشۡكُرُ لِنَفۡسِهِۦۖ وَمَن كَفَرَ فَإِنَّ رَبِّي غَنِيّٞ كَرِيمٞ
- [82:6](https://quran.com/82/6): يَٰٓأَيُّهَا ٱلۡإِنسَٰنُ مَا غَرَّكَ بِرَبِّكَ ٱلۡكَرِيمِ

### 43. Ar Raqeeb — الرَّقِيبُ

**Short label:** The Watchful One

**Main description:** Is the One who observes and beholds all that He has created (and the actions of creatures are also creations).

**Numbered-file description (unused):** Ar Raqeeb is one of the names of Allah in Islam, meaning 'The Observer' or 'The Watchful'.

**Displayed summary:** The name Ar Raqeeb signifies Allah's attribute of being ever-watchful and aware of all actions and deeds of His creation. This name emphasizes the concept of divine surveillance and accountability, illustrating that nothing escapes Allah's notice. Its origin is rooted in the Quran and Hadith, reflecting the belief that Allah observes all things with clarity and fairness.

**Audit:** Watchfulness is aligned; overlap with Shaheed/Hafeez is legitimate.

- [4:1](https://quran.com/4/1): يَٰٓأَيُّهَا ٱلنَّاسُ ٱتَّقُواْ رَبَّكُمُ ٱلَّذِي خَلَقَكُم مِّن نَّفۡسٖ وَٰحِدَةٖ وَخَلَقَ مِنۡهَا زَوۡجَهَا وَبَثَّ مِنۡهُمَا رِجَالٗا كَثِيرٗا وَنِسَآءٗۚ وَٱتَّقُواْ ٱللَّهَ ٱلَّذِي تَسَآءَلُونَ بِهِۦ وَٱلۡأَرۡحَامَۚ إِنَّ ٱللَّهَ كَانَ عَلَيۡكُمۡ رَقِيبٗا
- [5:117](https://quran.com/5/117): مَا قُلۡتُ لَهُمۡ إِلَّا مَآ أَمَرۡتَنِي بِهِۦٓ أَنِ ٱعۡبُدُواْ ٱللَّهَ رَبِّي وَرَبَّكُمۡۚ وَكُنتُ عَلَيۡهِمۡ شَهِيدٗا مَّا دُمۡتُ فِيهِمۡۖ فَلَمَّا تَوَفَّيۡتَنِي كُنتَ أَنتَ ٱلرَّقِيبَ عَلَيۡهِمۡۚ وَأَنتَ عَلَىٰ كُلِّ شَيۡءٖ شَهِيدٌ

### 44. Al Mujeeb — الْمُجِيبُ

**Short label:** The Responder to Prayer

**Main description:** Is the One who answers the invocations of His creatures if He wills and to whom He wills.

**Numbered-file description (unused):** Al Mujeeb is one of the names of Allah in Islam, meaning 'The Responsive One' or 'The Answerer of Prayers'.

**Displayed summary:** The name Al Mujeeb expresses Allah's attribute of responding to the supplications and needs of His creation. It emphasizes His willingness to answer prayers and grant requests, highlighting His mercy and compassion. The origin of this name is derived from Arabic, specifically from the root 'J-B', which relates to responding or answering. This name encourages believers to turn to Allah in their times of need, trusting in His responsiveness.

**Audit:** Trim trailing whitespace in transliteration. Correct incomplete J-B root to j-w-b (ج و ب). The Responsive is concise.

- [11:61](https://quran.com/11/61): ۞ وَإِلَىٰ ثَمُودَ أَخَاهُمۡ صَٰلِحٗاۚ قَالَ يَٰقَوۡمِ ٱعۡبُدُواْ ٱللَّهَ مَا لَكُم مِّنۡ إِلَٰهٍ غَيۡرُهُۥۖ هُوَ أَنشَأَكُم مِّنَ ٱلۡأَرۡضِ وَٱسۡتَعۡمَرَكُمۡ فِيهَا فَٱسۡتَغۡفِرُوهُ ثُمَّ تُوبُوٓاْ إِلَيۡهِۚ إِنَّ رَبِّي قَرِيبٞ مُّجِيبٞ

### 45. Al Waasi' — الْوَاسِعُ

**Short label:** The All-Encompassing

**Main description:** This name in Arabic means: immensely wide, but to qualify God it is understood that He has filled every void and has contained all things in a way that only He knows.

**Numbered-file description (unused):** Al Waasi' is one of the names of Allah in Islam, meaning 'The All-Encompassing' or 'The Boundless'.

**Displayed summary:** This name reflects Allah's infinite nature, His ability to encompass all things within His knowledge, mercy, and power. It signifies that Allah's provisions and knowledge extend beyond human comprehension and limitations, offering reassurance of His omnipresence and support in the lives of His creations.

**Audit:** Remove fills every void/contains all spatial claims and omnipresence wording. Explain encompassing knowledge, mercy and bounty.

- [2:268](https://quran.com/2/268): ٱلشَّيۡطَٰنُ يَعِدُكُمُ ٱلۡفَقۡرَ وَيَأۡمُرُكُم بِٱلۡفَحۡشَآءِۖ وَٱللَّهُ يَعِدُكُم مَّغۡفِرَةٗ مِّنۡهُ وَفَضۡلٗاۗ وَٱللَّهُ وَٰسِعٌ عَلِيمٞ
- [3:73](https://quran.com/3/73): وَلَا تُؤۡمِنُوٓاْ إِلَّا لِمَن تَبِعَ دِينَكُمۡ قُلۡ إِنَّ ٱلۡهُدَىٰ هُدَى ٱللَّهِ أَن يُؤۡتَىٰٓ أَحَدٞ مِّثۡلَ مَآ أُوتِيتُمۡ أَوۡ يُحَآجُّوكُمۡ عِندَ رَبِّكُمۡۗ قُلۡ إِنَّ ٱلۡفَضۡلَ بِيَدِ ٱللَّهِ يُؤۡتِيهِ مَن يَشَآءُۗ وَٱللَّهُ وَٰسِعٌ عَلِيمٞ
- [5:54](https://quran.com/5/54): يَٰٓأَيُّهَا ٱلَّذِينَ ءَامَنُواْ مَن يَرۡتَدَّ مِنكُمۡ عَن دِينِهِۦ فَسَوۡفَ يَأۡتِي ٱللَّهُ بِقَوۡمٖ يُحِبُّهُمۡ وَيُحِبُّونَهُۥٓ أَذِلَّةٍ عَلَى ٱلۡمُؤۡمِنِينَ أَعِزَّةٍ عَلَى ٱلۡكَٰفِرِينَ يُجَٰهِدُونَ فِي سَبِيلِ ٱللَّهِ وَلَا يَخَافُونَ لَوۡمَةَ لَآئِمٖۚ ذَٰلِكَ فَضۡلُ ٱللَّهِ يُؤۡتِيهِ مَن يَشَآءُۚ وَٱللَّهُ وَٰسِعٌ عَلِيمٌ

### 46. Al Hakeem — الْحَكِيمُ

**Short label:** The Perfectly Wise

**Main description:** He whose Wisdom has dazzled, dazzled and will dazzle any conscious and intelligent mind.

**Numbered-file description (unused):** Al Hakeem is one of the names of Allah in Islam, meaning 'The All-Wise.' It reflects Allah's perfect wisdom and knowledge in all matters.

**Displayed summary:** The name Al Hakeem emphasizes Allah's infinite understanding and the intricate wisdom behind His creation and decrees. It suggests that every action of Allah is grounded in wisdom, and believers are encouraged to trust in His plans even when they do not understand.

**Audit:** Replace repeated dazzled prose with wisdom in judgment, action and decree.

- [31:27](https://quran.com/31/27): وَلَوۡ أَنَّمَا فِي ٱلۡأَرۡضِ مِن شَجَرَةٍ أَقۡلَٰمٞ وَٱلۡبَحۡرُ يَمُدُّهُۥ مِنۢ بَعۡدِهِۦ سَبۡعَةُ أَبۡحُرٖ مَّا نَفِدَتۡ كَلِمَٰتُ ٱللَّهِۚ إِنَّ ٱللَّهَ عَزِيزٌ حَكِيمٞ
- [46:2](https://quran.com/46/2): تَنزِيلُ ٱلۡكِتَٰبِ مِنَ ٱللَّهِ ٱلۡعَزِيزِ ٱلۡحَكِيمِ
- [57:1](https://quran.com/57/1): سَبَّحَ لِلَّهِ مَا فِي ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِۖ وَهُوَ ٱلۡعَزِيزُ ٱلۡحَكِيمُ
- [66:2](https://quran.com/66/2): قَدۡ فَرَضَ ٱللَّهُ لَكُمۡ تَحِلَّةَ أَيۡمَٰنِكُمۡۚ وَٱللَّهُ مَوۡلَىٰكُمۡۖ وَهُوَ ٱلۡعَلِيمُ ٱلۡحَكِيمُ

### 47. Al Wadood — الْوَدُودُ

**Short label:** The Most Loving

**Main description:** This word in Arabic means: The one who draws near to his beloved by mercy, pity and beneficence, but when one reads this name in the "Quran", one feels each time different meanings which are difficult to translate into words. We can try to translate it: The one who draws near to His servants out of pity although they move away from Him and do not glorify Him.

**Numbered-file description (unused):** Al Wadood is one of the names of Allah, meaning 'The Most Loving' or 'The Affectionate.' This name signifies God's immense love and compassion towards His creation.

**Displayed summary:** Al Wadood reflects the quality of love that Allah has for His creations, showcasing the divine nature of compassion and affection. The name is derived from the Arabic root word 'Wad,' which means love or affection. It is mentioned in Islamic texts to emphasize God's nurturing and caring attributes.

**Audit:** Replace remarks about translation difficulty with a direct explanation of divine love; preserve corrected Wadood spelling.

- [11:90](https://quran.com/11/90): وَٱسۡتَغۡفِرُواْ رَبَّكُمۡ ثُمَّ تُوبُوٓاْ إِلَيۡهِۚ إِنَّ رَبِّي رَحِيمٞ وَدُودٞ
- [85:14](https://quran.com/85/14): وَهُوَ ٱلۡغَفُورُ ٱلۡوَدُودُ

### 48. Al Majeed — الْمَجِيدُ

**Short label:** The Glorious One, The Most Honorable

**Main description:** His glory encompasses nobility, excellence in deeds and generosity. Al Majeed and Al Maajid share this meaning, with a difference of emphasis.

**Numbered-file description (unused):** Al Majeed is one of the names of Allah in Islam, which translates to 'The Glorious' or 'The Most Majestic'. It signifies the perfection, greatness, and honor of Allah, highlighting His supreme status and divine attributes.

**Displayed summary:** The name Al Majeed emphasizes Allah's majesty and glory. It is derived from the Arabic root 'm-j-d', which relates to honor and greatness. This name reminds believers of the importance of recognizing and respecting the grandeur of Allah.

**Audit:** Use The All-Glorious or The Glorious; the current heading is six words. Do not invent an absolute distinction from Maajid.

- [11:73](https://quran.com/11/73): قَالُوٓاْ أَتَعۡجَبِينَ مِنۡ أَمۡرِ ٱللَّهِۖ رَحۡمَتُ ٱللَّهِ وَبَرَكَٰتُهُۥ عَلَيۡكُمۡ أَهۡلَ ٱلۡبَيۡتِۚ إِنَّهُۥ حَمِيدٞ مَّجِيدٞ

### 49. Al Baa'ith — الْبَاعِثُ

**Short label:** The Resurrector

**Main description:** Is the One who raises the dead and deploys life in whatever He wants, be they men, animals or simply earth.

**Numbered-file description (unused):** Al Baa'ith is one of the names of Allah in Islam, meaning 'The Resurrector' or 'The Resurrecting One'.

**Displayed summary:** The origin of the name Al Baa'ith comes from the Arabic root word 'baa'atha', which means 'to raise or to send forth'. In the context of Islamic theology, it refers to Allah's power to resurrect the dead and to bring about new life. This attribute is a reminder of the Day of Judgment, where all individuals will be brought back to life for accountability. It highlights Allah’s authority over life and death and serves to instill hope in believers for resurrection and eternal life.

**Audit:** 22:7 describes resurrection through a verb; label it a related attribute, not a literal occurrence of Al Baa'ith.

- [22:7](https://quran.com/22/7): وَأَنَّ ٱلسَّاعَةَ ءَاتِيَةٞ لَّا رَيۡبَ فِيهَا وَأَنَّ ٱللَّهَ يَبۡعَثُ مَن فِي ٱلۡقُبُورِ

### 50. Ash Shaheed — الشَّهِيدُ

**Short label:** The Witness

**Main description:** The Witness: Allah is the ultimate and perfect witness to all that has happened in His kingdom, with clear and blatant evidence, perfectly blameless.

**Numbered-file description (unused):** Ash Shaheed is one of the names of Allah, often translated as 'The All-Witnessing' or 'The Martyr'. It signifies the divine attribute of witnessing all actions and events, affirming the ultimate truth of reality, and reflects on God's justice and knowledge.

**Displayed summary:** The name Ash Shaheed represents God's omnipresence and omniscience, emphasizing that nothing escapes His knowledge. The origin of this name is rooted in Arabic, with 'Shaheed' deriving from the verb 'shahada', meaning to witness. It is a reminder for believers to act with integrity and righteousness, understanding that their actions are observed and recorded.

**Audit:** Witness is aligned; remove the unnecessary omnipresence claim from the explanation.

- [4:166](https://quran.com/4/166): لَّٰكِنِ ٱللَّهُ يَشۡهَدُ بِمَآ أَنزَلَ إِلَيۡكَۖ أَنزَلَهُۥ بِعِلۡمِهِۦۖ وَٱلۡمَلَٰٓئِكَةُ يَشۡهَدُونَۚ وَكَفَىٰ بِٱللَّهِ شَهِيدًا
- [22:17](https://quran.com/22/17): إِنَّ ٱلَّذِينَ ءَامَنُواْ وَٱلَّذِينَ هَادُواْ وَٱلصَّٰبِـِٔينَ وَٱلنَّصَٰرَىٰ وَٱلۡمَجُوسَ وَٱلَّذِينَ أَشۡرَكُوٓاْ إِنَّ ٱللَّهَ يَفۡصِلُ بَيۡنَهُمۡ يَوۡمَ ٱلۡقِيَٰمَةِۚ إِنَّ ٱللَّهَ عَلَىٰ كُلِّ شَيۡءٖ شَهِيدٌ
- [41:53](https://quran.com/41/53): سَنُرِيهِمۡ ءَايَٰتِنَا فِي ٱلۡأٓفَاقِ وَفِيٓ أَنفُسِهِمۡ حَتَّىٰ يَتَبَيَّنَ لَهُمۡ أَنَّهُ ٱلۡحَقُّۗ أَوَلَمۡ يَكۡفِ بِرَبِّكَ أَنَّهُۥ عَلَىٰ كُلِّ شَيۡءٖ شَهِيدٌ
- [48:28](https://quran.com/48/28): هُوَ ٱلَّذِيٓ أَرۡسَلَ رَسُولَهُۥ بِٱلۡهُدَىٰ وَدِينِ ٱلۡحَقِّ لِيُظۡهِرَهُۥ عَلَى ٱلدِّينِ كُلِّهِۦۚ وَكَفَىٰ بِٱللَّهِ شَهِيدٗا

### 51. Al Haqq — الْحَقُّ

**Short label:** The Truth

**Main description:** Is the One who holds the absolute Truth.

**Numbered-file description (unused):** Al Haqq, meaning 'The Absolute Truth,' is one of the names of Allah in Islam. It signifies the nature of Allah as the ultimate truth and reality, contrasting with falsehood and impermanence.

**Displayed summary:** The name Al Haqq highlights Allah's role as the source of all true knowledge and the ultimate arbiter of justice. It emphasizes that all matters will return to Him and that He is the foundation upon which truth and justice are built.

**Audit:** Truth/reality is aligned.

- [6:62](https://quran.com/6/62): ثُمَّ رُدُّوٓاْ إِلَى ٱللَّهِ مَوۡلَىٰهُمُ ٱلۡحَقِّۚ أَلَا لَهُ ٱلۡحُكۡمُ وَهُوَ أَسۡرَعُ ٱلۡحَٰسِبِينَ
- [22:6](https://quran.com/22/6): ذَٰلِكَ بِأَنَّ ٱللَّهَ هُوَ ٱلۡحَقُّ وَأَنَّهُۥ يُحۡيِ ٱلۡمَوۡتَىٰ وَأَنَّهُۥ عَلَىٰ كُلِّ شَيۡءٖ قَدِيرٞ
- [23:116](https://quran.com/23/116): فَتَعَٰلَى ٱللَّهُ ٱلۡمَلِكُ ٱلۡحَقُّۖ لَآ إِلَٰهَ إِلَّا هُوَ رَبُّ ٱلۡعَرۡشِ ٱلۡكَرِيمِ
- [24:25](https://quran.com/24/25): يَوۡمَئِذٖ يُوَفِّيهِمُ ٱللَّهُ دِينَهُمُ ٱلۡحَقَّ وَيَعۡلَمُونَ أَنَّ ٱللَّهَ هُوَ ٱلۡحَقُّ ٱلۡمُبِينُ

### 52. Al Wakeel — الْوَكِيلُ

**Short label:** The Trustee

**Main description:** Is the One who guarantees the sustenance of creatures and the One from whom nothing of their every state can be concealed

**Numbered-file description (unused):** Al Wakeel is one of the names of Allah in Islam, meaning 'The Trustee' or 'The Disposer of Affairs'. It signifies that Allah is the one who manages, protects, and oversees all affairs in the universe.

**Displayed summary:** The name implies that Allah is the ultimate protector and guardian of His creation, handling all matters and ensuring they are done justly and in accordance with His wisdom. The name Al Wakeel is derived from the Arabic root 'w-k-l', which means 'to appoint, to designate, or to entrust'. It appears in Islamic texts and emphasizes the reliance of creation on Allah's management.

**Audit:** Trustee/reliance is aligned; do not imply that reliance removes human responsibility.

- [3:173](https://quran.com/3/173): ٱلَّذِينَ قَالَ لَهُمُ ٱلنَّاسُ إِنَّ ٱلنَّاسَ قَدۡ جَمَعُواْ لَكُمۡ فَٱخۡشَوۡهُمۡ فَزَادَهُمۡ إِيمَٰنٗا وَقَالُواْ حَسۡبُنَا ٱللَّهُ وَنِعۡمَ ٱلۡوَكِيلُ
- [4:171](https://quran.com/4/171): يَٰٓأَهۡلَ ٱلۡكِتَٰبِ لَا تَغۡلُواْ فِي دِينِكُمۡ وَلَا تَقُولُواْ عَلَى ٱللَّهِ إِلَّا ٱلۡحَقَّۚ إِنَّمَا ٱلۡمَسِيحُ عِيسَى ٱبۡنُ مَرۡيَمَ رَسُولُ ٱللَّهِ وَكَلِمَتُهُۥٓ أَلۡقَىٰهَآ إِلَىٰ مَرۡيَمَ وَرُوحٞ مِّنۡهُۖ فَـَٔامِنُواْ بِٱللَّهِ وَرُسُلِهِۦۖ وَلَا تَقُولُواْ ثَلَٰثَةٌۚ ٱنتَهُواْ خَيۡرٗا لَّكُمۡۚ إِنَّمَا ٱللَّهُ إِلَٰهٞ وَٰحِدٞۖ سُبۡحَٰنَهُۥٓ أَن يَكُونَ لَهُۥ وَلَدٞۘ لَّهُۥ مَا فِي ٱلسَّمَٰوَٰتِ وَمَا فِي ٱلۡأَرۡضِۗ وَكَفَىٰ بِٱللَّهِ وَكِيلٗا
- [28:28](https://quran.com/28/28): قَالَ ذَٰلِكَ بَيۡنِي وَبَيۡنَكَۖ أَيَّمَا ٱلۡأَجَلَيۡنِ قَضَيۡتُ فَلَا عُدۡوَٰنَ عَلَيَّۖ وَٱللَّهُ عَلَىٰ مَا نَقُولُ وَكِيلٞ
- [73:9](https://quran.com/73/9): رَّبُّ ٱلۡمَشۡرِقِ وَٱلۡمَغۡرِبِ لَآ إِلَٰهَ إِلَّا هُوَ فَٱتَّخِذۡهُ وَكِيلٗا

### 53. Al Qawiyy — الْقَوِيُّ

**Short label:** The Possessor of all strength

**Main description:** The Mighty: Allah is the most powerful, the strongest, his strength has no limit and is not described with definitions of physical strength.

**Numbered-file description (unused):** Al Qawiyy is one of the names of Allah in Islam, meaning 'The All-Strong' or 'The Powerful'. It reflects God's supreme strength and ability over all creation.

**Displayed summary:** The name Al Qawiyy emphasizes the unmatched strength and power of Allah, highlighting that He is capable of doing all things without any limitations. This name originates from the Arabic root 'q-w-y', which signifies strength and might. Al Qawiyy reminds believers of the importance of relying on Allah's strength during times of difficulty.

**Audit:** The All-Strong is concise; explain complete strength without bodily comparisons.

- [22:40](https://quran.com/22/40): ٱلَّذِينَ أُخۡرِجُواْ مِن دِيَٰرِهِم بِغَيۡرِ حَقٍّ إِلَّآ أَن يَقُولُواْ رَبُّنَا ٱللَّهُۗ وَلَوۡلَا دَفۡعُ ٱللَّهِ ٱلنَّاسَ بَعۡضَهُم بِبَعۡضٖ لَّهُدِّمَتۡ صَوَٰمِعُ وَبِيَعٞ وَصَلَوَٰتٞ وَمَسَٰجِدُ يُذۡكَرُ فِيهَا ٱسۡمُ ٱللَّهِ كَثِيرٗاۗ وَلَيَنصُرَنَّ ٱللَّهُ مَن يَنصُرُهُۥٓۚ إِنَّ ٱللَّهَ لَقَوِيٌّ عَزِيزٌ
- [22:74](https://quran.com/22/74): مَا قَدَرُواْ ٱللَّهَ حَقَّ قَدۡرِهِۦٓۚ إِنَّ ٱللَّهَ لَقَوِيٌّ عَزِيزٌ
- [42:19](https://quran.com/42/19): ٱللَّهُ لَطِيفُۢ بِعِبَادِهِۦ يَرۡزُقُ مَن يَشَآءُۖ وَهُوَ ٱلۡقَوِيُّ ٱلۡعَزِيزُ
- [57:25](https://quran.com/57/25): لَقَدۡ أَرۡسَلۡنَا رُسُلَنَا بِٱلۡبَيِّنَٰتِ وَأَنزَلۡنَا مَعَهُمُ ٱلۡكِتَٰبَ وَٱلۡمِيزَانَ لِيَقُومَ ٱلنَّاسُ بِٱلۡقِسۡطِۖ وَأَنزَلۡنَا ٱلۡحَدِيدَ فِيهِ بَأۡسٞ شَدِيدٞ وَمَنَٰفِعُ لِلنَّاسِ وَلِيَعۡلَمَ ٱللَّهُ مَن يَنصُرُهُۥ وَرُسُلَهُۥ بِٱلۡغَيۡبِۚ إِنَّ ٱللَّهَ قَوِيٌّ عَزِيزٞ

### 54. Al Mateen — الْمَتِينُ

**Short label:** The Forceful

**Main description:** The Enduring: An endurance worthy of Him, unlimited, indefatigable, and which is not marred by any weakness, without pause, without interruption and which no factor can attenuate.

**Numbered-file description (unused):** Al Mateen is one of the names of Allah in Islam, meaning 'The Firm' or 'The Steadfast.' It signifies Allah's power, strength, and unwavering support for His creation.

**Displayed summary:** Al Mateen highlights the strength and reliability of Allah, showcasing that He is a source of firmness in faith and an unyielding support for believers. The name emphasizes that Allah is ever-present and capable of providing strength to those who seek Him.

**Audit:** The Firm is clearer than The Forceful; strength is unfailing, not merely forceful conduct.

- [51:58](https://quran.com/51/58): إِنَّ ٱللَّهَ هُوَ ٱلرَّزَّاقُ ذُو ٱلۡقُوَّةِ ٱلۡمَتِينُ

### 55. Al Waliyy — الْوَلِيُّ

**Short label:** The Protector

**Main description:** Is the One who preserves the believing servants, so the prophets and those who follow them are those who are the "preserved people" in the sense that they have been guided to the obligatory thanks to Him who makes them exist, that is, that is to say with regard to Allah

**Numbered-file description (unused):** Al Waliyy is one of the names of Allah in Islam. It translates to 'The Protecting Associate' or 'The Patron.' This name signifies Allah's role as a protector and guide for His creation, emphasizing His closeness and support to those who believe in Him.

**Displayed summary:** The name Al Waliyy reflects Allah's protective nature and His intimate relationship with His believers. Originating from Arabic, it highlights the concept of divine guardianship, where Allah assists, nurtures, and supports those who seek His help, ensuring their well-being in both spiritual and worldly matters.

**Audit:** Replace rambling main description with protecting and supporting believers. Waliyy and Waali are distinct forms with overlapping meanings.

- [4:45](https://quran.com/4/45): وَٱللَّهُ أَعۡلَمُ بِأَعۡدَآئِكُمۡۚ وَكَفَىٰ بِٱللَّهِ وَلِيّٗا وَكَفَىٰ بِٱللَّهِ نَصِيرٗا
- [7:196](https://quran.com/7/196): إِنَّ وَلِـِّۧيَ ٱللَّهُ ٱلَّذِي نَزَّلَ ٱلۡكِتَٰبَۖ وَهُوَ يَتَوَلَّى ٱلصَّٰلِحِينَ
- [42:28](https://quran.com/42/28): وَهُوَ ٱلَّذِي يُنَزِّلُ ٱلۡغَيۡثَ مِنۢ بَعۡدِ مَا قَنَطُواْ وَيَنشُرُ رَحۡمَتَهُۥۚ وَهُوَ ٱلۡوَلِيُّ ٱلۡحَمِيدُ
- [45:19](https://quran.com/45/19): إِنَّهُمۡ لَن يُغۡنُواْ عَنكَ مِنَ ٱللَّهِ شَيۡـٔٗاۚ وَإِنَّ ٱلظَّٰلِمِينَ بَعۡضُهُمۡ أَوۡلِيَآءُ بَعۡضٖۖ وَٱللَّهُ وَلِيُّ ٱلۡمُتَّقِينَ

### 56. Al Hameed — الْحَمِيدُ

**Short label:** The Praised

**Main description:** Is the One who deserves the thanksgiving, the glorification and the song of praise more than any other

**Numbered-file description (unused):** Al Hameed is one of the names of Allah in Islam, meaning 'The Praiseworthy.' It denotes Allah's perfection, worthiness of praise, and the excellence of His attributes.

**Displayed summary:** The name Al Hameed emphasizes the significance of Allah's nature as deserving of all praise and commendation. It recognizes that all good qualities and actions originate from Him, and humanity is encouraged to acknowledge and express gratitude for His blessings. The term comes from the Arabic root 'ḥ-m-d,' which relates to praise and commendation.

**Audit:** Praise is aligned; The Praiseworthy is more idiomatic than The Praised if changing the label.

- [14:8](https://quran.com/14/8): وَقَالَ مُوسَىٰٓ إِن تَكۡفُرُوٓاْ أَنتُمۡ وَمَن فِي ٱلۡأَرۡضِ جَمِيعٗا فَإِنَّ ٱللَّهَ لَغَنِيٌّ حَمِيدٌ
- [31:12](https://quran.com/31/12): وَلَقَدۡ ءَاتَيۡنَا لُقۡمَٰنَ ٱلۡحِكۡمَةَ أَنِ ٱشۡكُرۡ لِلَّهِۚ وَمَن يَشۡكُرۡ فَإِنَّمَا يَشۡكُرُ لِنَفۡسِهِۦۖ وَمَن كَفَرَ فَإِنَّ ٱللَّهَ غَنِيٌّ حَمِيدٞ
- [31:26](https://quran.com/31/26): لِلَّهِ مَا فِي ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِۚ إِنَّ ٱللَّهَ هُوَ ٱلۡغَنِيُّ ٱلۡحَمِيدُ
- [41:42](https://quran.com/41/42): لَّا يَأۡتِيهِ ٱلۡبَٰطِلُ مِنۢ بَيۡنِ يَدَيۡهِ وَلَا مِنۡ خَلۡفِهِۦۖ تَنزِيلٞ مِّنۡ حَكِيمٍ حَمِيدٖ

### 57. Al Muhsi — الْمُحْصِي

**Short label:** The All-Enumerating

**Main description:** Literally translated as the counter but means the One who counts every creature, every act, every speck of dust in his kingdom and of whom nothing escapes his knowledge

**Numbered-file description (unused):** Al Muhsi is one of the names of Allah in Islam, meanings 'The All-Enumerating' or 'The Counter'.

**Displayed summary:** The name Al Muhsi signifies Allah’s attribute of being fully aware of everything that exists, counting all actions and deeds, ensuring nothing is overlooked. This name emphasizes the concept of accountability in the hereafter, where every action will be recorded and judged. The origin of this name is from the Arabic root 'حَصَرَ' which means to enumerate or count.

**Audit:** Correct root حصر to ح ص ي. 72:28 and 78:29 use verbs of enumeration. 82:10–12 describes recording angels and also creates an unsupported range URL.

- [72:28](https://quran.com/72/28): لِّيَعۡلَمَ أَن قَدۡ أَبۡلَغُواْ رِسَٰلَٰتِ رَبِّهِمۡ وَأَحَاطَ بِمَا لَدَيۡهِمۡ وَأَحۡصَىٰ كُلَّ شَيۡءٍ عَدَدَۢا
- [78:29](https://quran.com/78/29): وَكُلَّ شَيۡءٍ أَحۡصَيۡنَٰهُ كِتَٰبٗا
- [82:10](https://quran.com/82/10): وَإِنَّ عَلَيۡكُمۡ لَحَٰفِظِينَ
- [82:11](https://quran.com/82/11): كِرَامٗا كَٰتِبِينَ
- [82:12](https://quran.com/82/12): يَعۡلَمُونَ مَا تَفۡعَلُونَ

### 58. Al Mubdi — الْمُبْدِئُ

**Short label:** The Originator

**Main description:** Is the One who gives a beginning to what does not exist without a beginning, without a model, that is to say, the One who creates creatures

**Numbered-file description (unused):** Al Mubdi, meaning 'The Originator', is one of the names of Allah in Islam. It signifies the act of creating something from nothing and initiating existence.

**Displayed summary:** The name Al Mubdi represents Allah's attribute of being the Creator and the one who brings forth life and creation. It emphasizes the divine power to start all things and highlights the concept of renewal and new beginnings.

**Audit:** Origination is aligned. All existing references use verbs, not the listed noun.

- [10:34](https://quran.com/10/34): قُلۡ هَلۡ مِن شُرَكَآئِكُم مَّن يَبۡدَؤُاْ ٱلۡخَلۡقَ ثُمَّ يُعِيدُهُۥۚ قُلِ ٱللَّهُ يَبۡدَؤُاْ ٱلۡخَلۡقَ ثُمَّ يُعِيدُهُۥۖ فَأَنَّىٰ تُؤۡفَكُونَ
- [27:64](https://quran.com/27/64): أَمَّن يَبۡدَؤُاْ ٱلۡخَلۡقَ ثُمَّ يُعِيدُهُۥ وَمَن يَرۡزُقُكُم مِّنَ ٱلسَّمَآءِ وَٱلۡأَرۡضِۗ أَءِلَٰهٞ مَّعَ ٱللَّهِۚ قُلۡ هَاتُواْ بُرۡهَٰنَكُمۡ إِن كُنتُمۡ صَٰدِقِينَ
- [29:19](https://quran.com/29/19): أَوَلَمۡ يَرَوۡاْ كَيۡفَ يُبۡدِئُ ٱللَّهُ ٱلۡخَلۡقَ ثُمَّ يُعِيدُهُۥٓۚ إِنَّ ذَٰلِكَ عَلَى ٱللَّهِ يَسِيرٞ
- [85:13](https://quran.com/85/13): إِنَّهُۥ هُوَ يُبۡدِئُ وَيُعِيدُ

### 59. Al Mu'eed — الْمُعِيدُ

**Short label:** The Restorer

**Main description:** Is the One who brings the living back to death after life and the One who brings the dead back to life after death

**Numbered-file description (unused):** Al Mu'eed is one of the names of Allah in Islam. It means 'The Restorer' or 'The Reinstater'.

**Displayed summary:** The origin of the name Al Mu'eed comes from the Arabic root word 'ع و د' (ʿ-w-d), which signifies the act of returning or bringing back. In Islamic belief, this name emphasizes Allah's attribute of restoring creation and granting second chances, whether in terms of guidance, forgiveness, or the revival of life after death.

**Audit:** Correct the description's direction: restoring creation/resurrection, not bringing the living back to death. Remove the summary's unsupported focus on second chances/forgiveness. References use verbs.

- [10:34](https://quran.com/10/34): قُلۡ هَلۡ مِن شُرَكَآئِكُم مَّن يَبۡدَؤُاْ ٱلۡخَلۡقَ ثُمَّ يُعِيدُهُۥۚ قُلِ ٱللَّهُ يَبۡدَؤُاْ ٱلۡخَلۡقَ ثُمَّ يُعِيدُهُۥۖ فَأَنَّىٰ تُؤۡفَكُونَ
- [27:64](https://quran.com/27/64): أَمَّن يَبۡدَؤُاْ ٱلۡخَلۡقَ ثُمَّ يُعِيدُهُۥ وَمَن يَرۡزُقُكُم مِّنَ ٱلسَّمَآءِ وَٱلۡأَرۡضِۗ أَءِلَٰهٞ مَّعَ ٱللَّهِۚ قُلۡ هَاتُواْ بُرۡهَٰنَكُمۡ إِن كُنتُمۡ صَٰدِقِينَ
- [29:19](https://quran.com/29/19): أَوَلَمۡ يَرَوۡاْ كَيۡفَ يُبۡدِئُ ٱللَّهُ ٱلۡخَلۡقَ ثُمَّ يُعِيدُهُۥٓۚ إِنَّ ذَٰلِكَ عَلَى ٱللَّهِ يَسِيرٞ
- [85:13](https://quran.com/85/13): إِنَّهُۥ هُوَ يُبۡدِئُ وَيُعِيدُ

### 60. Al Muhyee — الْمُحْيِي

**Short label:** The Giver of life

**Main description:** Is the One who gives life to a mixture of maniyy, that is to say of seminal liquid, devoid of soul and makes him a living being and the One who gives life to the bodies of humans, jinns and angels by reuniting them with their souls during the resurrection

**Numbered-file description (unused):** Al Muhyee is one of the names of Allah in Islam, meaning 'The Giver of Life' or 'The Restorer.' It refers to Allah's power to grant life and revitalize creation.

**Displayed summary:** The name Al Muhyee originates from the Arabic root 'حيا' (ḥ-y-y), which is associated with life and living. This name emphasizes the belief that Allah is the source of all life and has the ability to give it and restore it. It reflects His mercy and control over the cycles of life and death.

**Audit:** 30:50 has the nominal form muhyi; other references use verbs. Remove unnecessary biological and resurrection details from the description.

- [7:158](https://quran.com/7/158): قُلۡ يَٰٓأَيُّهَا ٱلنَّاسُ إِنِّي رَسُولُ ٱللَّهِ إِلَيۡكُمۡ جَمِيعًا ٱلَّذِي لَهُۥ مُلۡكُ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِۖ لَآ إِلَٰهَ إِلَّا هُوَ يُحۡيِۦ وَيُمِيتُۖ فَـَٔامِنُواْ بِٱللَّهِ وَرَسُولِهِ ٱلنَّبِيِّ ٱلۡأُمِّيِّ ٱلَّذِي يُؤۡمِنُ بِٱللَّهِ وَكَلِمَٰتِهِۦ وَٱتَّبِعُوهُ لَعَلَّكُمۡ تَهۡتَدُونَ
- [15:23](https://quran.com/15/23): وَإِنَّا لَنَحۡنُ نُحۡيِۦ وَنُمِيتُ وَنَحۡنُ ٱلۡوَٰرِثُونَ
- [30:50](https://quran.com/30/50): فَٱنظُرۡ إِلَىٰٓ ءَاثَٰرِ رَحۡمَتِ ٱللَّهِ كَيۡفَ يُحۡيِ ٱلۡأَرۡضَ بَعۡدَ مَوۡتِهَآۚ إِنَّ ذَٰلِكَ لَمُحۡيِ ٱلۡمَوۡتَىٰۖ وَهُوَ عَلَىٰ كُلِّ شَيۡءٖ قَدِيرٞ
- [57:2](https://quran.com/57/2): لَهُۥ مُلۡكُ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِۖ يُحۡيِۦ وَيُمِيتُۖ وَهُوَ عَلَىٰ كُلِّ شَيۡءٖ قَدِيرٌ

### 61. Al Mumeet — اَلْمُمِيتُ

**Short label:** The Taker of life

**Main description:** He is the One who kills living beings and the One who cancels out the power of the powerful of this world by death

**Numbered-file description (unused):** Al Mumeet, meaning 'The Creator of Death', is one of the names of Allah in Islam. It signifies Allah's authority over life and death.

**Displayed summary:** The name Al Mumeet originates from the Arabic root word 'moot', which means death. In Islamic belief, it emphasizes that Allah is the one who decrees death and has complete control over the end of life for all living beings. This attribute encourages reflection on mortality and inspires believers to lead a purposeful life.

**Audit:** Life-taking is aligned; all listed references use verbs rather than the noun Al Mumeet.

- [3:156](https://quran.com/3/156): يَٰٓأَيُّهَا ٱلَّذِينَ ءَامَنُواْ لَا تَكُونُواْ كَٱلَّذِينَ كَفَرُواْ وَقَالُواْ لِإِخۡوَٰنِهِمۡ إِذَا ضَرَبُواْ فِي ٱلۡأَرۡضِ أَوۡ كَانُواْ غُزّٗى لَّوۡ كَانُواْ عِندَنَا مَا مَاتُواْ وَمَا قُتِلُواْ لِيَجۡعَلَ ٱللَّهُ ذَٰلِكَ حَسۡرَةٗ فِي قُلُوبِهِمۡۗ وَٱللَّهُ يُحۡيِۦ وَيُمِيتُۗ وَٱللَّهُ بِمَا تَعۡمَلُونَ بَصِيرٞ
- [7:158](https://quran.com/7/158): قُلۡ يَٰٓأَيُّهَا ٱلنَّاسُ إِنِّي رَسُولُ ٱللَّهِ إِلَيۡكُمۡ جَمِيعًا ٱلَّذِي لَهُۥ مُلۡكُ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِۖ لَآ إِلَٰهَ إِلَّا هُوَ يُحۡيِۦ وَيُمِيتُۖ فَـَٔامِنُواْ بِٱللَّهِ وَرَسُولِهِ ٱلنَّبِيِّ ٱلۡأُمِّيِّ ٱلَّذِي يُؤۡمِنُ بِٱللَّهِ وَكَلِمَٰتِهِۦ وَٱتَّبِعُوهُ لَعَلَّكُمۡ تَهۡتَدُونَ
- [15:23](https://quran.com/15/23): وَإِنَّا لَنَحۡنُ نُحۡيِۦ وَنُمِيتُ وَنَحۡنُ ٱلۡوَٰرِثُونَ
- [57:2](https://quran.com/57/2): لَهُۥ مُلۡكُ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِۖ يُحۡيِۦ وَيُمِيتُۖ وَهُوَ عَلَىٰ كُلِّ شَيۡءٖ قَدِيرٌ

### 62. Al Haiyy — الْحَيُّ

**Short label:** The Ever Living

**Main description:** The Living One: Is the One who does not come into existence and yet exists, and the One who possesses a hayah and this attribute is not a life as is the hayah of creatures made up of a combination of soul, flesh and of blood Imam At-Tahawiyy said in his treatise on belief that "he who attributes to Allah an attribute in a sense which is among those which can be attributed to human beings is disbeliever"

**Numbered-file description (unused):** Al Haiyy translates to "The Ever-Living" in Arabic. It is one of the names of Allah in Islam, reflecting His eternal existence and the fact that He is alive without beginning or end.

**Displayed summary:** The name Al Haiyy signifies the attribute of life that is all-encompassing and perfect. It conveys that Allah's life is not like that of His creation; He is self-sustained and independent, contrasting with the finite nature of all living beings. This name emphasizes the importance of Allah's eternal presence and His continuous maintenance of the universe.

**Audit:** Remove the long unrelated quotation about disbelief from the description. Explain life without death, sleep or dependence.

- [2:255](https://quran.com/2/255): ٱللَّهُ لَآ إِلَٰهَ إِلَّا هُوَ ٱلۡحَيُّ ٱلۡقَيُّومُۚ لَا تَأۡخُذُهُۥ سِنَةٞ وَلَا نَوۡمٞۚ لَّهُۥ مَا فِي ٱلسَّمَٰوَٰتِ وَمَا فِي ٱلۡأَرۡضِۗ مَن ذَا ٱلَّذِي يَشۡفَعُ عِندَهُۥٓ إِلَّا بِإِذۡنِهِۦۚ يَعۡلَمُ مَا بَيۡنَ أَيۡدِيهِمۡ وَمَا خَلۡفَهُمۡۖ وَلَا يُحِيطُونَ بِشَيۡءٖ مِّنۡ عِلۡمِهِۦٓ إِلَّا بِمَا شَآءَۚ وَسِعَ كُرۡسِيُّهُ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضَۖ وَلَا يَـُٔودُهُۥ حِفۡظُهُمَاۚ وَهُوَ ٱلۡعَلِيُّ ٱلۡعَظِيمُ
- [3:2](https://quran.com/3/2): ٱللَّهُ لَآ إِلَٰهَ إِلَّا هُوَ ٱلۡحَيُّ ٱلۡقَيُّومُ
- [25:58](https://quran.com/25/58): وَتَوَكَّلۡ عَلَى ٱلۡحَيِّ ٱلَّذِي لَا يَمُوتُ وَسَبِّحۡ بِحَمۡدِهِۦۚ وَكَفَىٰ بِهِۦ بِذُنُوبِ عِبَادِهِۦ خَبِيرًا
- [40:65](https://quran.com/40/65): هُوَ ٱلۡحَيُّ لَآ إِلَٰهَ إِلَّا هُوَ فَٱدۡعُوهُ مُخۡلِصِينَ لَهُ ٱلدِّينَۗ ٱلۡحَمۡدُ لِلَّهِ رَبِّ ٱلۡعَٰلَمِينَ

### 63. Al Qayyoom — الْقَيُّومُ

**Short label:** The Self Existing

**Main description:** Is the One who exists, is eternal and does not change and the One who is not affected by non-existence or annihilation and is in no way affected by controlling and creating the acts of human beings

**Numbered-file description (unused):** Al Qayyoom is one of the names of Allah in Islam, and it translates to 'The Sustainer' or 'The Holder of All Things.' It signifies Allah's eternal existence and His role in maintaining and sustaining the universe.

**Displayed summary:** The name Al Qayyoom emphasizes Allah's self-sufficiency and His absolute power in sustaining His creation. It highlights the concept that everything depends on Him for existence and continuity, as He is always present and aware of all that happens. This name is derived from the Arabic root 'Qama,' which means 'to stand' or 'to maintain.' Al Qayyoom reflects the idea that Allah is the everlasting and ever-watchful guardian of all beings.

**Audit:** Description must include sustaining all existence as well as self-subsistence; summary already includes this. The Sustainer is concise but may legitimately overlap other labels.

- [2:255](https://quran.com/2/255): ٱللَّهُ لَآ إِلَٰهَ إِلَّا هُوَ ٱلۡحَيُّ ٱلۡقَيُّومُۚ لَا تَأۡخُذُهُۥ سِنَةٞ وَلَا نَوۡمٞۚ لَّهُۥ مَا فِي ٱلسَّمَٰوَٰتِ وَمَا فِي ٱلۡأَرۡضِۗ مَن ذَا ٱلَّذِي يَشۡفَعُ عِندَهُۥٓ إِلَّا بِإِذۡنِهِۦۚ يَعۡلَمُ مَا بَيۡنَ أَيۡدِيهِمۡ وَمَا خَلۡفَهُمۡۖ وَلَا يُحِيطُونَ بِشَيۡءٖ مِّنۡ عِلۡمِهِۦٓ إِلَّا بِمَا شَآءَۚ وَسِعَ كُرۡسِيُّهُ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضَۖ وَلَا يَـُٔودُهُۥ حِفۡظُهُمَاۚ وَهُوَ ٱلۡعَلِيُّ ٱلۡعَظِيمُ
- [3:2](https://quran.com/3/2): ٱللَّهُ لَآ إِلَٰهَ إِلَّا هُوَ ٱلۡحَيُّ ٱلۡقَيُّومُ
- [20:111](https://quran.com/20/111): ۞ وَعَنَتِ ٱلۡوُجُوهُ لِلۡحَيِّ ٱلۡقَيُّومِۖ وَقَدۡ خَابَ مَنۡ حَمَلَ ظُلۡمٗا

### 64. Al Waajid — الْوَاجِدُ

**Short label:** The Perceiver

**Main description:** Is the One who does not need anything and is not reached or concerned with want, lack or need

**Numbered-file description (unused):** Al Waajid is one of the names of Allah in Islam, meaning 'The Perceiver' or 'The Finder'.

**Displayed summary:** Al Waajid signifies Allah's ability to perceive and identify all things, conveying that nothing is hidden from Him. This name reflects the belief that Allah knows everything and can find what may be lost to human perception. It originates from Arabic, with 'waajid' derived from the root that encompasses concepts of presence and realization.

**Audit:** Substantive mismatch: Perceiver/summary of perception versus description of lacking nothing. Resolve using a named commentary before choosing a label. 38:44 says Allah found Ayyub patient; it does not call Allah Al Waajid.

- [38:44](https://quran.com/38/44): وَخُذۡ بِيَدِكَ ضِغۡثٗا فَٱضۡرِب بِّهِۦ وَلَا تَحۡنَثۡۗ إِنَّا وَجَدۡنَٰهُ صَابِرٗاۚ نِّعۡمَ ٱلۡعَبۡدُ إِنَّهُۥٓ أَوَّابٞ

### 65. Al Maajid — الْمَاجِدُ

**Short label:** The Glorious

**Main description:** This name shares with Al Majeed the meaning of glory, nobility and generosity; it does not denote an unrelated attribute.

**Numbered-file description (unused):** Al Maajid is one of the names of Allah in Islam, translating to 'The Glorious' or 'The Majestic.' It signifies the greatness of Allah's attributes and His exalted nature.

**Displayed summary:** The name Al Maajid emphasizes the glory and majesty of Allah. It highlights His honor and the reverence He commands in the hearts of believers. The origin of this name comes from the Arabic root related to honor and magnificence, reflecting the Islamic belief in Allah's supremacy and grandeur in all aspects of existence.

**Audit:** Glory is aligned; 11:73 contains Majeed, not Maajid. Mark as a related form and retain the distinction in Arabic without inventing a rigid English distinction.

- [11:73](https://quran.com/11/73): قَالُوٓاْ أَتَعۡجَبِينَ مِنۡ أَمۡرِ ٱللَّهِۖ رَحۡمَتُ ٱللَّهِ وَبَرَكَٰتُهُۥ عَلَيۡكُمۡ أَهۡلَ ٱلۡبَيۡتِۚ إِنَّهُۥ حَمِيدٞ مَّجِيدٞ

### 66. Al Waahid — الْواحِدُ

**Short label:** The Only One

**Main description:** The Unique: Is the One who has no alter ego, whether in His Being, in His Act, or in His Attributes. Nothing is like him

**Numbered-file description (unused):** Al Waahid means 'The One' in Arabic. This name signifies the oneness and uniqueness of Allah, emphasizing that there is no partner or equal to Him.

**Displayed summary:** Al Waahid reflects the concept of divine unity in Islam, underscoring the belief that Allah is singular and incomparable. This name is derived from the Arabic root 'W-H-D', which conveys the idea of oneness and individuality. It is a reminder of Allah's sole sovereignty and His singular nature, reinforcing the central tenet of Tawhid in Islamic belief.

**Audit:** Oneness is aligned; Waahid and Ahad are different words expressing overlapping truths, not accidental duplicates.

- [2:163](https://quran.com/2/163): وَإِلَٰهُكُمۡ إِلَٰهٞ وَٰحِدٞۖ لَّآ إِلَٰهَ إِلَّا هُوَ ٱلرَّحۡمَٰنُ ٱلرَّحِيمُ
- [5:73](https://quran.com/5/73): لَّقَدۡ كَفَرَ ٱلَّذِينَ قَالُوٓاْ إِنَّ ٱللَّهَ ثَالِثُ ثَلَٰثَةٖۘ وَمَا مِنۡ إِلَٰهٍ إِلَّآ إِلَٰهٞ وَٰحِدٞۚ وَإِن لَّمۡ يَنتَهُواْ عَمَّا يَقُولُونَ لَيَمَسَّنَّ ٱلَّذِينَ كَفَرُواْ مِنۡهُمۡ عَذَابٌ أَلِيمٌ
- [9:31](https://quran.com/9/31): ٱتَّخَذُوٓاْ أَحۡبَارَهُمۡ وَرُهۡبَٰنَهُمۡ أَرۡبَابٗا مِّن دُونِ ٱللَّهِ وَٱلۡمَسِيحَ ٱبۡنَ مَرۡيَمَ وَمَآ أُمِرُوٓاْ إِلَّا لِيَعۡبُدُوٓاْ إِلَٰهٗا وَٰحِدٗاۖ لَّآ إِلَٰهَ إِلَّا هُوَۚ سُبۡحَٰنَهُۥ عَمَّا يُشۡرِكُونَ
- [18:110](https://quran.com/18/110): قُلۡ إِنَّمَآ أَنَا۠ بَشَرٞ مِّثۡلُكُمۡ يُوحَىٰٓ إِلَيَّ أَنَّمَآ إِلَٰهُكُمۡ إِلَٰهٞ وَٰحِدٞۖ فَمَن كَانَ يَرۡجُواْ لِقَآءَ رَبِّهِۦ فَلۡيَعۡمَلۡ عَمَلٗا صَٰلِحٗا وَلَا يُشۡرِكۡ بِعِبَادَةِ رَبِّهِۦٓ أَحَدَۢا

### 67. Al Ahad — اَلاَحَدُ

**Short label:** The One

**Main description:** The One: Is He Who has uniqueness as an attribute

**Numbered-file description (unused):** Al Ahad is one of the names of Allah in Islam, signifying 'The One' or 'The Unique'. This name emphasizes the oneness and indivisibility of God, underscoring that He has no partners, equals, or rivals.

**Displayed summary:** The name Al Ahad is derived from the Arabic root word 'ahd', which means a covenant or promise. In Islamic theology, it represents the singularity of God's essence and nature. Al Ahad denotes that Allah is the sole creator and sustainer of the universe, reinforcing the concept of monotheism (Tawhid) in Islam.

**Audit:** Remove the false covenant etymology: أحد (Ahad) is not عهد (ahd). Rebuild summary around unique oneness and 112:1–4.

- [112:1](https://quran.com/112/1): قُلۡ هُوَ ٱللَّهُ أَحَدٌ

### 68. As Samad — الصَّمَدُ

**Short label:** The Eternal Refuge

**Main description:** Allah is wholly self-sufficient, needing nothing. All creation depends on Him and turns to Him for the fulfilment of its needs.

**Numbered-file description (unused):** As Samad is one of the names of Allah in Islam, which translates to 'The Eternal' or 'The Absolute.' It refers to Allah's self-sufficiency and the idea that He is in need of nothing and that everything relies on Him.

**Displayed summary:** The name As Samad originates from the Quran, specifically from Surah Al-Ikhlas (Chapter 112), which emphasizes the unique oneness and eternal nature of Allah. It is a reminder of God's independence and the source of all sustenance and existence.

**Audit:** Keep agreed The Eternal Refuge and revised main description. Update separate summary to state both needing nothing and all creation turning to Him in need; current summary is too generic.

- [112:2](https://quran.com/112/2): ٱللَّهُ ٱلصَّمَدُ

### 69. Al Qaadir — الْقَادِرُ

**Short label:** The Powerful

**Main description:** Is the One who is unaffected by weakness, the One whose power is absolute and dominates all creatures and the One whom absolutely nothing can weaken

**Numbered-file description (unused):** Al Qaadir is one of the names of Allah in Islam, meaning 'The Powerful' or 'The Omnipotent'. It signifies Allah's absolute power and ability to do all things.

**Displayed summary:** The name Al Qaadir is derived from the Arabic root 'q-d-r', which denotes strength and capability. This name reflects the belief that Allah has complete control over the universe and is capable of exercising His will in all matters. It emphasizes the concept of divine power and authority in Islamic theology.

**Audit:** Power is aligned; preserve meaningful overlap with Muqtadir.

- [6:65](https://quran.com/6/65): قُلۡ هُوَ ٱلۡقَادِرُ عَلَىٰٓ أَن يَبۡعَثَ عَلَيۡكُمۡ عَذَابٗا مِّن فَوۡقِكُمۡ أَوۡ مِن تَحۡتِ أَرۡجُلِكُمۡ أَوۡ يَلۡبِسَكُمۡ شِيَعٗا وَيُذِيقَ بَعۡضَكُم بَأۡسَ بَعۡضٍۗ ٱنظُرۡ كَيۡفَ نُصَرِّفُ ٱلۡأٓيَٰتِ لَعَلَّهُمۡ يَفۡقَهُونَ
- [36:81](https://quran.com/36/81): أَوَلَيۡسَ ٱلَّذِي خَلَقَ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضَ بِقَٰدِرٍ عَلَىٰٓ أَن يَخۡلُقَ مِثۡلَهُمۚ بَلَىٰ وَهُوَ ٱلۡخَلَّٰقُ ٱلۡعَلِيمُ
- [46:33](https://quran.com/46/33): أَوَلَمۡ يَرَوۡاْ أَنَّ ٱللَّهَ ٱلَّذِي خَلَقَ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضَ وَلَمۡ يَعۡيَ بِخَلۡقِهِنَّ بِقَٰدِرٍ عَلَىٰٓ أَن يُحۡـِۧيَ ٱلۡمَوۡتَىٰۚ بَلَىٰٓۚ إِنَّهُۥ عَلَىٰ كُلِّ شَيۡءٖ قَدِيرٞ
- [75:40](https://quran.com/75/40): أَلَيۡسَ ذَٰلِكَ بِقَٰدِرٍ عَلَىٰٓ أَن يُحۡـِۧيَ ٱلۡمَوۡتَىٰ

### 70. Al Muqtadir — الْمُقْتَدِرُ

**Short label:** The Creator of all power

**Main description:** Is the One who is mighty over all things and the One whose capacities are absolutely limitless

**Numbered-file description (unused):** Al Muqtadir is one of the names of Allah in Islam, which means 'The Powerful' or 'The Omnipotent.' This name conveys the concept of absolute power and authority, indicating that Allah has complete control over all things.

**Displayed summary:** The name signifies Allah's capability to execute His will and manage all affairs of creation. It emphasizes His strength and the fact that nothing occurs without His permission. Derived from the Arabic root 'q-d-r,' which relates to power and capability, Al Muqtadir is mentioned in the Quran, reflecting the attributes of strength and supreme authority.

**Audit:** Prefer The All-Powerful over Creator of all power; describe complete power and capability.

- [18:45](https://quran.com/18/45): وَٱضۡرِبۡ لَهُم مَّثَلَ ٱلۡحَيَوٰةِ ٱلدُّنۡيَا كَمَآءٍ أَنزَلۡنَٰهُ مِنَ ٱلسَّمَآءِ فَٱخۡتَلَطَ بِهِۦ نَبَاتُ ٱلۡأَرۡضِ فَأَصۡبَحَ هَشِيمٗا تَذۡرُوهُ ٱلرِّيَٰحُۗ وَكَانَ ٱللَّهُ عَلَىٰ كُلِّ شَيۡءٖ مُّقۡتَدِرًا
- [54:42](https://quran.com/54/42): كَذَّبُواْ بِـَٔايَٰتِنَا كُلِّهَا فَأَخَذۡنَٰهُمۡ أَخۡذَ عَزِيزٖ مُّقۡتَدِرٍ
- [54:55](https://quran.com/54/55): فِي مَقۡعَدِ صِدۡقٍ عِندَ مَلِيكٖ مُّقۡتَدِرِۭ

### 71. Al Muqaddim — الْمُقَدِّمُ

**Short label:** The Expediter

**Main description:** Is He who is free from the characteristics of creatures and better than any creature; He who fixes the speed of events, accelerates certain things in the eyes of people according to a destiny and grants each thing a value, a rank, a time and characteristics which suit Him according to wisdom

**Numbered-file description (unused):** One of the 99 names of Allah in Islam, meaning 'The Expediter' or 'The Promoter.'

**Displayed summary:** Al Muqaddim signifies Allah's ability to advance and prioritize matters as He wills. This name reflects His divine wisdom in the way He orders events in the universe, ensuring that everything unfolds according to His plan. It emphasizes the concept that God is in control of all situations, facilitating and hastening what needs to occur in a person's life or the world at large.

**Audit:** Remove unrelated 17:34. 16:61 concerns appointed terms, not literal Al Muqaddim. Muslim 771a contains the name explicitly in Arabic.

- [16:61](https://quran.com/16/61): وَلَوۡ يُؤَاخِذُ ٱللَّهُ ٱلنَّاسَ بِظُلۡمِهِم مَّا تَرَكَ عَلَيۡهَا مِن دَآبَّةٖ وَلَٰكِن يُؤَخِّرُهُمۡ إِلَىٰٓ أَجَلٖ مُّسَمّٗىۖ فَإِذَا جَآءَ أَجَلُهُمۡ لَا يَسۡتَـٔۡخِرُونَ سَاعَةٗ وَلَا يَسۡتَقۡدِمُونَ
- [17:34](https://quran.com/17/34): وَلَا تَقۡرَبُواْ مَالَ ٱلۡيَتِيمِ إِلَّا بِٱلَّتِي هِيَ أَحۡسَنُ حَتَّىٰ يَبۡلُغَ أَشُدَّهُۥۚ وَأَوۡفُواْ بِٱلۡعَهۡدِۖ إِنَّ ٱلۡعَهۡدَ كَانَ مَسۡـُٔولٗا

### 72. Al Mu’akhkhir — الْمُؤَخِّرُ

**Short label:** The Delayer

**Main description:** Is He who is free from the characteristics of creatures and better than any creature; The one who intends the slowness of events and slows down certain things in the eyes of people according to a destiny and grants each thing a value, a rank, a time and characteristics which suit it according to a wisdom (cf no 71)

**Numbered-file description (unused):** Al Mu'akhkhir is one of the names of Allah, which means 'The Delayer'. It refers to God's control over time and His ability to postpone events according to His wisdom.

**Displayed summary:** The name Al Mu'akhkhir emphasizes the concept that Allah has the ultimate authority over what happens when, reminding believers that everything unfolds in the timing that is best according to divine wisdom. Its origin can be traced to the Arabic language and Islamic teachings, reflecting the understanding of God's nature as encompassing all aspects of time and existence.

**Audit:** 71:4 uses a verb of delaying, not the listed noun. Use Muslim 771a for the explicit name; keep corrected transliteration.

- [71:4](https://quran.com/71/4): يَغۡفِرۡ لَكُم مِّن ذُنُوبِكُمۡ وَيُؤَخِّرۡكُمۡ إِلَىٰٓ أَجَلٖ مُّسَمًّىۚ إِنَّ أَجَلَ ٱللَّهِ إِذَا جَآءَ لَا يُؤَخَّرُۚ لَوۡ كُنتُمۡ تَعۡلَمُونَ

### 73. Al Awwal — الأوَّلُ

**Short label:** The First

**Main description:** The First: Is the One who exists without creatures, without beginning and without time; One who is exempt from being created, from the appearance of a new attribute and from change, One who exists without coming into existence

**Numbered-file description (unused):** Al Awwal is one of the names of Allah in Islam, meaning 'The First.' It signifies that Allah is the first and has no beginning, emphasizing His eternal and timeless nature.

**Displayed summary:** The name Al Awwal illustrates Allah's precedence over all creation, indicating that everything else comes after Him in existence. It serves to remind believers of His sovereignty and the notion that He transcends time and sequence.

**Audit:** 57:3 is direct. Use Muslim 2713a's explanation that nothing precedes Him.

- [57:3](https://quran.com/57/3): هُوَ ٱلۡأَوَّلُ وَٱلۡأٓخِرُ وَٱلظَّٰهِرُ وَٱلۡبَاطِنُۖ وَهُوَ بِكُلِّ شَيۡءٍ عَلِيمٌ

### 74. Al Aakhir — الآخِرُ

**Short label:** The Last

**Main description:** The Last: Is the One who exists without interruption while the creatures are annihilated at every moment; He who is eternal, He who annihilates and is exempt from being annihilated and the disappearance of one or more attributes; The one that exists without being concerned by the notion of end and annihilation

**Numbered-file description (unused):** Al Aakhir is one of the names of Allah in Islam, meaning 'The Last'.

**Displayed summary:** The name Al Aakhir signifies that Allah is the ultimate end of all existence. In Islamic theology, this name emphasizes His eternal nature and the idea that everything in the universe has a beginning and an end, but Allah remains beyond the constraints of time. It reflects the belief in the Day of Judgment and the finality of His will.

**Audit:** 57:3 is direct. Use Muslim 2713a's explanation that nothing comes after Him.

- [57:3](https://quran.com/57/3): هُوَ ٱلۡأَوَّلُ وَٱلۡأٓخِرُ وَٱلظَّٰهِرُ وَٱلۡبَاطِنُۖ وَهُوَ بِكُلِّ شَيۡءٍ عَلِيمٌ

### 75. Az Zaahir — الظَّاهِرُ

**Short label:** The Manifest

**Main description:** Is the One whose existence is manifested in everything by proofs, from the point of view of power and domination and not what would be present in any body from the point of view of place, image and how

**Numbered-file description (unused):** Az-Zaahir is one of the names of Allah in Islam, meaning 'The Manifest' or 'The Evident.' It signifies the visible aspect of God's essence, reflecting His presence and the clarity of His attributes.

**Displayed summary:** The name Az-Zaahir illustrates Allah's nature as being apparent and known, emphasizing that His attributes and creations are clear and evident in the world. This name encourages believers to recognize the signs of God in the universe.

**Audit:** 57:3 is direct. Explain with Muslim 2713a; avoid silently choosing a contested spatial paraphrase.

- [57:3](https://quran.com/57/3): هُوَ ٱلۡأَوَّلُ وَٱلۡأٓخِرُ وَٱلظَّٰهِرُ وَٱلۡبَاطِنُۖ وَهُوَ بِكُلِّ شَيۡءٍ عَلِيمٌ

### 76. Al Baatin — الْبَاطِنُ

**Short label:** The Hidden

**Main description:** Is the One who is preserved from the delusions of creatures attempting to attribute to the Creator human characteristics or characteristics of creatures and the One who creates the characteristics of creatures and their images

**Numbered-file description (unused):** Al Baatin, meaning 'The Hidden One,' refers to one of the names of Allah in Islam. It denotes the aspect of God's nature that is not visible or apparent to human understanding. This name emphasizes that God's essence and attributes are beyond human perception and comprehension.

**Displayed summary:** Al Baatin signifies the hidden and unseen nature of Allah, reflecting the belief that while Allah is manifest in creation, His true essence remains concealed from human sight and understanding. This name reminds believers of the limitations of human knowledge and the depth of God's presence in the universe.

**Audit:** 57:3 is direct. Explain with Muslim 2713a; do not equate hiddenness with being absent or unknown in every sense.

- [57:3](https://quran.com/57/3): هُوَ ٱلۡأَوَّلُ وَٱلۡأٓخِرُ وَٱلظَّٰهِرُ وَٱلۡبَاطِنُۖ وَهُوَ بِكُلِّ شَيۡءٍ عَلِيمٌ

### 77. Al Waali — الْوَالِي

**Short label:** The Governor

**Main description:** Is the One who owns, manages and controls all things and the One who creates everything according to wisdom

**Numbered-file description (unused):** Al Waali is one of the names of Allah in Islamic tradition, meaning 'The Sole Governor' or 'The Patron'.

**Displayed summary:** The name Al Waali signifies Allah's role as the ultimate authority and protector over all creation. It indicates His sovereignty and ability to manage and govern the affairs of the universe. The origin of this name can be traced to Arabic roots, where it emphasizes the loving and guiding relationship Allah has with His creation, showcasing His mercy and power.

**Audit:** 13:11 has waalin, a corresponding indefinite form meaning protector. Remove unrelated 22:7, which is about resurrection.

- [13:11](https://quran.com/13/11): لَهُۥ مُعَقِّبَٰتٞ مِّنۢ بَيۡنِ يَدَيۡهِ وَمِنۡ خَلۡفِهِۦ يَحۡفَظُونَهُۥ مِنۡ أَمۡرِ ٱللَّهِۗ إِنَّ ٱللَّهَ لَا يُغَيِّرُ مَا بِقَوۡمٍ حَتَّىٰ يُغَيِّرُواْ مَا بِأَنفُسِهِمۡۗ وَإِذَآ أَرَادَ ٱللَّهُ بِقَوۡمٖ سُوٓءٗا فَلَا مَرَدَّ لَهُۥۚ وَمَا لَهُم مِّن دُونِهِۦ مِن وَالٍ
- [22:7](https://quran.com/22/7): وَأَنَّ ٱلسَّاعَةَ ءَاتِيَةٞ لَّا رَيۡبَ فِيهَا وَأَنَّ ٱللَّهَ يَبۡعَثُ مَن فِي ٱلۡقُبُورِ

### 78. Al Muta’ali — الْمُتَعَالِي

**Short label:** The Supreme One

**Main description:** Is the One who is free from the attributes of creatures and the One whose power over creatures is boundless

**Numbered-file description (unused):** Al Muta'ali (المتعالي) is one of the names of Allah in Islam, meaning 'The Transcendent' or 'The Most Exalted'. This name emphasizes Allah's absolute majesty and superiority over all of creation.

**Displayed summary:** Al Muta'ali refers to Allah's nature of being beyond all limitations and imperfections, highlighting His greatness and the fact that He is above anything that can be imagined or compared. The name inspires reverence and acknowledgment of Allah's unique status as the Creator, who is exalted above all that exists.

**Audit:** 13:9 directly contains al-muta'al. Supremacy/exaltedness is aligned.

- [13:9](https://quran.com/13/9): عَٰلِمُ ٱلۡغَيۡبِ وَٱلشَّهَٰدَةِ ٱلۡكَبِيرُ ٱلۡمُتَعَالِ

### 79. Al Barr — الْبَرُّ

**Short label:** The Doer of Good

**Main description:** Is the One who has the ability to bestow faith and blessings on all creatures, and among them some are grateful and some are not: the former are called Muslims or believers and others are called non-Muslims or non-Muslims. -believers

**Numbered-file description (unused):** Al-Barr is one of the names of Allah in Islam, which translates to 'The Source of All Goodness' or 'The Kind One'.

**Displayed summary:** The name Al-Barr signifies Allah's attributes of kindness, compassion, and benevolence. It conveys the idea that Allah is the ultimate source of all goodness and moral uprightness, encouraging believers to embody these qualities in their own lives. The name is derived from the Arabic root word 'B-R-R', which implies righteousness and devotion.

**Audit:** The Benevolent is concise. Remove duplicated non-Muslims text and repair the main description.

- [52:28](https://quran.com/52/28): إِنَّا كُنَّا مِن قَبۡلُ نَدۡعُوهُۖ إِنَّهُۥ هُوَ ٱلۡبَرُّ ٱلرَّحِيمُ

### 80. At Tawwaab — التَّوَابُ

**Short label:** The Guide to Repentance

**Main description:** Allah repeatedly enables His servants to repent and accepts their repentance.

**Numbered-file description (unused):** At Tawwaab is one of the names of Allah in Islam, meaning 'The Acceptor of Repentance'. It reflects God's mercy and willingness to forgive His servants who turn back to Him in sincere repentance.

**Displayed summary:** The name signifies Allah's infinite mercy and compassion, emphasizing His readiness to forgive those who seek forgiveness after sinning. At Tawwaab is derived from the Arabic root 't-w-b', which relates to the act of returning or repenting. It underscores the importance of repentance (tawbah) in Islam.

**Audit:** Change Guide to Repentance to Accepter of Repentance; align summary with both enabling repentance and accepting it, as in the revised description.

- [2:128](https://quran.com/2/128): رَبَّنَا وَٱجۡعَلۡنَا مُسۡلِمَيۡنِ لَكَ وَمِن ذُرِّيَّتِنَآ أُمَّةٗ مُّسۡلِمَةٗ لَّكَ وَأَرِنَا مَنَاسِكَنَا وَتُبۡ عَلَيۡنَآۖ إِنَّكَ أَنتَ ٱلتَّوَّابُ ٱلرَّحِيمُ
- [4:64](https://quran.com/4/64): وَمَآ أَرۡسَلۡنَا مِن رَّسُولٍ إِلَّا لِيُطَاعَ بِإِذۡنِ ٱللَّهِۚ وَلَوۡ أَنَّهُمۡ إِذ ظَّلَمُوٓاْ أَنفُسَهُمۡ جَآءُوكَ فَٱسۡتَغۡفَرُواْ ٱللَّهَ وَٱسۡتَغۡفَرَ لَهُمُ ٱلرَّسُولُ لَوَجَدُواْ ٱللَّهَ تَوَّابٗا رَّحِيمٗا
- [49:12](https://quran.com/49/12): يَٰٓأَيُّهَا ٱلَّذِينَ ءَامَنُواْ ٱجۡتَنِبُواْ كَثِيرٗا مِّنَ ٱلظَّنِّ إِنَّ بَعۡضَ ٱلظَّنِّ إِثۡمٞۖ وَلَا تَجَسَّسُواْ وَلَا يَغۡتَب بَّعۡضُكُم بَعۡضًاۚ أَيُحِبُّ أَحَدُكُمۡ أَن يَأۡكُلَ لَحۡمَ أَخِيهِ مَيۡتٗا فَكَرِهۡتُمُوهُۚ وَٱتَّقُواْ ٱللَّهَۚ إِنَّ ٱللَّهَ تَوَّابٞ رَّحِيمٞ
- [110:3](https://quran.com/110/3): فَسَبِّحۡ بِحَمۡدِ رَبِّكَ وَٱسۡتَغۡفِرۡهُۚ إِنَّهُۥ كَانَ تَوَّابَۢا

### 81. Al Muntaqim — الْمُنْتَقِمُ

**Short label:** The Avenger

**Main description:** Is the avenger

**Numbered-file description (unused):** Al Muntaqim is one of the names of Allah in Islam. It is derived from the Arabic root word 'naqama,' which means 'to take vengeance' or 'to retaliate.'

**Displayed summary:** Al Muntaqim, meaning 'The Avenger,' signifies Allah's attribute of holding accountable those who do wrong and seeking justice. This name reflects the divine nature of ensuring moral order and punishing transgressors in accordance with His wisdom and mercy. It serves as a reminder of the importance of justice in human conduct and the repercussions of one's actions.

**Audit:** Expand the three-word description to just requital for wrongdoing. References use plural muntaqimun, not singular Al Muntaqim; show that grammatical distinction.

- [32:22](https://quran.com/32/22): وَمَنۡ أَظۡلَمُ مِمَّن ذُكِّرَ بِـَٔايَٰتِ رَبِّهِۦ ثُمَّ أَعۡرَضَ عَنۡهَآۚ إِنَّا مِنَ ٱلۡمُجۡرِمِينَ مُنتَقِمُونَ
- [43:41](https://quran.com/43/41): فَإِمَّا نَذۡهَبَنَّ بِكَ فَإِنَّا مِنۡهُم مُّنتَقِمُونَ
- [44:16](https://quran.com/44/16): يَوۡمَ نَبۡطِشُ ٱلۡبَطۡشَةَ ٱلۡكُبۡرَىٰٓ إِنَّا مُنتَقِمُونَ

### 82. Al Afuww — العَفُوُّ

**Short label:** The Forgiver

**Main description:** Allah pardons by erasing sins. Al Ghafoor emphasises forgiveness and concealment of sins; Al Afuww emphasises their removal.

**Numbered-file description (unused):** Al Afuww is one of the names of Allah in Islam, meaning 'The Pardoner' or 'The Forgiver'. It highlights Allah's boundless mercy and forgiveness towards His creation, emphasizing His ability to absolve sins and grant pardon to those who repent sincerely.

**Displayed summary:** The concept of Al Afuww signifies Allah's readiness to forgive and His compassion towards His creatures. The name is derived from the Arabic root 'عَفَا' (ʿafā), which means to pardon or to overlook. This name is often invoked in prayers and supplications, reflecting the hope for divine mercy and forgiveness. It teaches believers the importance of repentance and the potential for redemption.

**Audit:** Change The Forgiver to The Pardoner, matching the revised description of effacing sins; forgiveness and pardon still overlap.

- [4:99](https://quran.com/4/99): فَأُوْلَٰٓئِكَ عَسَى ٱللَّهُ أَن يَعۡفُوَ عَنۡهُمۡۚ وَكَانَ ٱللَّهُ عَفُوًّا غَفُورٗا
- [4:149](https://quran.com/4/149): إِن تُبۡدُواْ خَيۡرًا أَوۡ تُخۡفُوهُ أَوۡ تَعۡفُواْ عَن سُوٓءٖ فَإِنَّ ٱللَّهَ كَانَ عَفُوّٗا قَدِيرًا
- [22:60](https://quran.com/22/60): ۞ ذَٰلِكَۖ وَمَنۡ عَاقَبَ بِمِثۡلِ مَا عُوقِبَ بِهِۦ ثُمَّ بُغِيَ عَلَيۡهِ لَيَنصُرَنَّهُ ٱللَّهُۚ إِنَّ ٱللَّهَ لَعَفُوٌّ غَفُورٞ

### 83. Ar Ra’oof — الرَّؤُوفُ

**Short label:** The Most Kind

**Main description:** Is the One who can offer an immensity of benefits without counterpart and without necessity

**Numbered-file description (unused):** Ar Ra'oof is one of the names of Allah in Islam, which translates to 'The Most Kind', 'The Most Compassionate', or 'The Most Caring'.

**Displayed summary:** The name Ar Ra'oof emphasizes Allah's deep compassion and kindness towards His creation. It reflects His tender care, particularly towards those who are in need of mercy and support. This attribute inspires believers to exhibit kindness and compassion in their own lives, encouraging them to emulate this divine quality in their interactions with others.

**Audit:** Kindness/compassion is aligned.

- [3:30](https://quran.com/3/30): يَوۡمَ تَجِدُ كُلُّ نَفۡسٖ مَّا عَمِلَتۡ مِنۡ خَيۡرٖ مُّحۡضَرٗا وَمَا عَمِلَتۡ مِن سُوٓءٖ تَوَدُّ لَوۡ أَنَّ بَيۡنَهَا وَبَيۡنَهُۥٓ أَمَدَۢا بَعِيدٗاۗ وَيُحَذِّرُكُمُ ٱللَّهُ نَفۡسَهُۥۗ وَٱللَّهُ رَءُوفُۢ بِٱلۡعِبَادِ
- [9:117](https://quran.com/9/117): لَّقَد تَّابَ ٱللَّهُ عَلَى ٱلنَّبِيِّ وَٱلۡمُهَٰجِرِينَ وَٱلۡأَنصَارِ ٱلَّذِينَ ٱتَّبَعُوهُ فِي سَاعَةِ ٱلۡعُسۡرَةِ مِنۢ بَعۡدِ مَا كَادَ يَزِيغُ قُلُوبُ فَرِيقٖ مِّنۡهُمۡ ثُمَّ تَابَ عَلَيۡهِمۡۚ إِنَّهُۥ بِهِمۡ رَءُوفٞ رَّحِيمٞ
- [57:9](https://quran.com/57/9): هُوَ ٱلَّذِي يُنَزِّلُ عَلَىٰ عَبۡدِهِۦٓ ءَايَٰتِۭ بَيِّنَٰتٖ لِّيُخۡرِجَكُم مِّنَ ٱلظُّلُمَٰتِ إِلَى ٱلنُّورِۚ وَإِنَّ ٱللَّهَ بِكُمۡ لَرَءُوفٞ رَّحِيمٞ
- [59:10](https://quran.com/59/10): وَٱلَّذِينَ جَآءُو مِنۢ بَعۡدِهِمۡ يَقُولُونَ رَبَّنَا ٱغۡفِرۡ لَنَا وَلِإِخۡوَٰنِنَا ٱلَّذِينَ سَبَقُونَا بِٱلۡإِيمَٰنِ وَلَا تَجۡعَلۡ فِي قُلُوبِنَا غِلّٗا لِّلَّذِينَ ءَامَنُواْ رَبَّنَآ إِنَّكَ رَءُوفٞ رَّحِيمٌ

### 84. Maalik Ul Mulk — مَالِكُ الْمُلْكِ

**Short label:** Master of the Kingdom

**Main description:** This name signifies Allah's absolute and complete control over all dominion, emphasizing that He is the sole possessor and ruler of everything that exists.

**Numbered-file description (unused):** Maalik Ul Mulk, which translates to 'Master of the Kingdom', is one of the names of Allah in Islam. It emphasizes His sovereignty and authority over all creation and dominions.

**Displayed summary:** The name 'Maalik Ul Mulk' is derived from Arabic, where 'Maalik' means 'Master' and 'Mulk' means 'Kingdom' or 'Sovereignty'. It signifies Allah's complete control over the universe, reinforcing the concept that all power and authority belong solely to Him. This name reflects His role as the ultimate ruler and caretaker of all that exists.

**Audit:** Sovereign of Dominion is a concise option; 3:26 contains the phrase directly.

- [3:26](https://quran.com/3/26): قُلِ ٱللَّهُمَّ مَٰلِكَ ٱلۡمُلۡكِ تُؤۡتِي ٱلۡمُلۡكَ مَن تَشَآءُ وَتَنزِعُ ٱلۡمُلۡكَ مِمَّن تَشَآءُ وَتُعِزُّ مَن تَشَآءُ وَتُذِلُّ مَن تَشَآءُۖ بِيَدِكَ ٱلۡخَيۡرُۖ إِنَّكَ عَلَىٰ كُلِّ شَيۡءٖ قَدِيرٞ

### 85. Dhu Al Jalaali Wa Al Ikraam — ذُوالْجَلاَلِ وَالإكْرَامِ

**Short label:** Possessor of Majesty and Bounty

**Main description:** That is, Allah is the One who is better than any other in reality and all the evidence confirms this and denies the contrary so it is not valid to dispute this point, to deny it or to be hostile against this subject; is the One who fills with blessings those who will have understood this who will know the general holiness, will have success and will be well guided on the Day of Judgment

**Numbered-file description (unused):** This name translates to 'The Lord of Glory and Honour' or 'The Lord of Majesty and Generosity'. It emphasizes Allah's attributes of majesty, greatness, and generosity towards His creation.

**Displayed summary:** Dhu Al Jalaali Wa Al Ikraam is one of the names of Allah in Islam, found in Hadith literature and referred to in the context of Allah's overwhelming majesty and benevolence. The origin of this name reflects the Islamic understanding of God's dual nature of being exalted with glory and simultaneously generous, reminding believers of His magnificence and the kindness He bestows upon them.

**Audit:** Remove God's dual nature. Name is explicitly in 55:27 and 55:78, not merely hadith. Explain majesty and generosity/honour. Do not force a three-word label if it loses half the compound meaning.

- [55:27](https://quran.com/55/27): وَيَبۡقَىٰ وَجۡهُ رَبِّكَ ذُو ٱلۡجَلَٰلِ وَٱلۡإِكۡرَامِ
- [55:78](https://quran.com/55/78): تَبَٰرَكَ ٱسۡمُ رَبِّكَ ذِي ٱلۡجَلَٰلِ وَٱلۡإِكۡرَامِ

### 86. Al Muqsit — الْمُقْسِطُ

**Short label:** The Equitable One

**Main description:** Is the One who judges according to the greatest justice, the One who is free from all forms of injustice and tyranny, and the One who has absolutely no accountability

**Numbered-file description (unused):** Al Muqsit is one of the names of Allah in Islam, meaning 'The Just' or 'The Equitable'. It reflects His attribute of fairness and justice in all matters.

**Displayed summary:** Al Muqsit highlights the concept of justice which is central to Islamic teachings. It emphasizes that Allah is fair and equitable in His judgments and dealings with creation. The origin of this name is rooted in the Quran and Hadith, where Allah's attributes are described, and it serves as a reminder for believers to strive for justice in their own lives.

**Audit:** 3:18 describes upholding justice; 7:29 commands justice. Neither literally calls Allah Al Muqsit. Mark as related attributes.

- [7:29](https://quran.com/7/29): قُلۡ أَمَرَ رَبِّي بِٱلۡقِسۡطِۖ وَأَقِيمُواْ وُجُوهَكُمۡ عِندَ كُلِّ مَسۡجِدٖ وَٱدۡعُوهُ مُخۡلِصِينَ لَهُ ٱلدِّينَۚ كَمَا بَدَأَكُمۡ تَعُودُونَ
- [3:18](https://quran.com/3/18): شَهِدَ ٱللَّهُ أَنَّهُۥ لَآ إِلَٰهَ إِلَّا هُوَ وَٱلۡمَلَٰٓئِكَةُ وَأُوْلُواْ ٱلۡعِلۡمِ قَآئِمَۢا بِٱلۡقِسۡطِۚ لَآ إِلَٰهَ إِلَّا هُوَ ٱلۡعَزِيزُ ٱلۡحَكِيمُ

### 87. Al Jaami' — الْجَامِعُ

**Short label:** The Gatherer

**Main description:** Is the One who gathers creatures together on a day of which there is no doubt, the Day of Judgment

**Numbered-file description (unused):** Al Jaami' is one of the names of Allah in Islam, meaning 'The Gatherer' or 'The Collector'. It signifies the divine attribute of bringing together all of creation, including souls on the Day of Judgment.

**Displayed summary:** The name Al Jaami' reflects God's power to gather all beings and souls, emphasizing His ultimate control over the universe and the resurrection of humanity. This concept is rooted in Islamic teachings about the Day of Judgment, where all individuals will be resurrected and held accountable for their actions.

**Audit:** 3:9 has jami', a direct nominal form referring to gathering people. Meaning is aligned.

- [3:9](https://quran.com/3/9): رَبَّنَآ إِنَّكَ جَامِعُ ٱلنَّاسِ لِيَوۡمٖ لَّا رَيۡبَ فِيهِۚ إِنَّ ٱللَّهَ لَا يُخۡلِفُ ٱلۡمِيعَادَ

### 88. Al Ghaniyy — الْغَنِيُّ

**Short label:** The Rich One

**Main description:** Is the One who has no need of creatures and whose creatures cannot do without al-Ghaniyy [who is the One who makes them exist and keeps them in existence and constantly changes them]

**Numbered-file description (unused):** Al Ghaniyy is one of the names of Allah in Islam, meaning 'The Self-Sufficient' or 'The Wealthy.' It represents Allah's attribute of being free of need, independent, and possessing all richness and sufficiency.

**Displayed summary:** The name Al Ghaniyy emphasizes Allah's self-sufficiency and independence from any external support. It reassures believers that Allah is the ultimate provider, in contrast to human dependency and need. The concept originates from the Quran, where Allah is described as possessing all wealth and sufficiency, highlighting His greatness and the assurance to His creation.

**Audit:** Prefer The Self-Sufficient to The Rich One, avoiding a merely financial reading. Explain no dependence on creation.

- [3:97](https://quran.com/3/97): فِيهِ ءَايَٰتُۢ بَيِّنَٰتٞ مَّقَامُ إِبۡرَٰهِيمَۖ وَمَن دَخَلَهُۥ كَانَ ءَامِنٗاۗ وَلِلَّهِ عَلَى ٱلنَّاسِ حِجُّ ٱلۡبَيۡتِ مَنِ ٱسۡتَطَاعَ إِلَيۡهِ سَبِيلٗاۚ وَمَن كَفَرَ فَإِنَّ ٱللَّهَ غَنِيٌّ عَنِ ٱلۡعَٰلَمِينَ
- [39:7](https://quran.com/39/7): إِن تَكۡفُرُواْ فَإِنَّ ٱللَّهَ غَنِيٌّ عَنكُمۡۖ وَلَا يَرۡضَىٰ لِعِبَادِهِ ٱلۡكُفۡرَۖ وَإِن تَشۡكُرُواْ يَرۡضَهُ لَكُمۡۗ وَلَا تَزِرُ وَازِرَةٞ وِزۡرَ أُخۡرَىٰۚ ثُمَّ إِلَىٰ رَبِّكُم مَّرۡجِعُكُمۡ فَيُنَبِّئُكُم بِمَا كُنتُمۡ تَعۡمَلُونَۚ إِنَّهُۥ عَلِيمُۢ بِذَاتِ ٱلصُّدُورِ
- [47:38](https://quran.com/47/38): هَٰٓأَنتُمۡ هَٰٓؤُلَآءِ تُدۡعَوۡنَ لِتُنفِقُواْ فِي سَبِيلِ ٱللَّهِ فَمِنكُم مَّن يَبۡخَلُۖ وَمَن يَبۡخَلۡ فَإِنَّمَا يَبۡخَلُ عَن نَّفۡسِهِۦۚ وَٱللَّهُ ٱلۡغَنِيُّ وَأَنتُمُ ٱلۡفُقَرَآءُۚ وَإِن تَتَوَلَّوۡاْ يَسۡتَبۡدِلۡ قَوۡمًا غَيۡرَكُمۡ ثُمَّ لَا يَكُونُوٓاْ أَمۡثَٰلَكُم
- [57:24](https://quran.com/57/24): ٱلَّذِينَ يَبۡخَلُونَ وَيَأۡمُرُونَ ٱلنَّاسَ بِٱلۡبُخۡلِۗ وَمَن يَتَوَلَّ فَإِنَّ ٱللَّهَ هُوَ ٱلۡغَنِيُّ ٱلۡحَمِيدُ

### 89. Al Mughni — الْمُغْنِي

**Short label:** The Enricher

**Main description:** Is the One who satisfies the needs of creatures and the One who brings them their sustenance

**Numbered-file description (unused):** Al Mughni is one of the names of Allah in Islam, which translates to 'The Inflictor of Help' or 'The One Who Grants Victory'.

**Displayed summary:** The name Al Mughni signifies Allah's attribute of providing support and assistance to those in need. It represents the divine capacity to offer help during times of struggle, guiding believers towards triumph over adversity. The origin of this name is rooted in Islamic teachings and is mentioned in various sources, emphasizing Allah's role as a protector and helper.

**Audit:** Replace separate description Inflictor of Help/Grants Victory and summary about triumph. Enriching/giving sufficiency is the relevant meaning. 9:28 and suggested 53:48 use verbs.

- [9:28](https://quran.com/9/28): يَٰٓأَيُّهَا ٱلَّذِينَ ءَامَنُوٓاْ إِنَّمَا ٱلۡمُشۡرِكُونَ نَجَسٞ فَلَا يَقۡرَبُواْ ٱلۡمَسۡجِدَ ٱلۡحَرَامَ بَعۡدَ عَامِهِمۡ هَٰذَاۚ وَإِنۡ خِفۡتُمۡ عَيۡلَةٗ فَسَوۡفَ يُغۡنِيكُمُ ٱللَّهُ مِن فَضۡلِهِۦٓ إِن شَآءَۚ إِنَّ ٱللَّهَ عَلِيمٌ حَكِيمٞ

### 90. Al Maani' — اَلْمَانِعُ

**Short label:** The Preventer of harm

**Main description:** Is the One who intends for each creature his victories, his protections and the One from whom they must always be preserved according to a destiny

**Numbered-file description (unused):** Al Maani' is one of the names of Allah in Islam, meaning 'The Withholder' or 'The Preventer'.

**Displayed summary:** This name reflects the divine attribute of Allah who controls and restrains that which is harmful or unnecessary for His creation. It highlights the concept that Allah has the authority to grant or withhold benefits and blessings, emphasizing His wisdom and sovereignty over all matters.

**Audit:** Change Preventer of harm to The Withholder. 67:21 concerns withholding provision; it does not promise universal protection from harm or contain this name as a noun.

- [67:21](https://quran.com/67/21): أَمَّنۡ هَٰذَا ٱلَّذِي يَرۡزُقُكُمۡ إِنۡ أَمۡسَكَ رِزۡقَهُۥۚ بَل لَّجُّواْ فِي عُتُوّٖ وَنُفُورٍ

### 91. Ad Daaarr — الضَّارَّ

**Short label:** The Creator of the harmful

**Main description:** Is the One who is almighty to bring inevitably the nuisance which must reach a creature according to a destiny

**Numbered-file description (unused):** Ad-Daaarr (Arabic: الضار) is one of the names of Allah in Islam, meaning 'The Distresser' or 'The Creator of Distress'.

**Displayed summary:** Ad-Daaarr reflects Allah's capability to create hardship or distress as part of His divine will. This name signifies that trials and difficulties can be a means of testing faith and character. Originating from the Arabic root ‘Daal-Ra’-‘Ayn, it emphasizes Allah's power over all aspects of existence, including the adverse circumstances one may face. It serves as a reminder that even in hardship, there is a greater wisdom and purpose at play.

**Audit:** Correct Daaarr to Daarr and false root Daal-Ra-Ayn to ض ر ر. Present hardship in relation to benefit and divine wisdom, with no implication of injustice. 6:17 is related, not a literal name occurrence.

- [6:17](https://quran.com/6/17): وَإِن يَمۡسَسۡكَ ٱللَّهُ بِضُرّٖ فَلَا كَاشِفَ لَهُۥٓ إِلَّا هُوَۖ وَإِن يَمۡسَسۡكَ بِخَيۡرٖ فَهُوَ عَلَىٰ كُلِّ شَيۡءٖ قَدِيرٞ

### 92. An Naafi’ — النَّافِعُ

**Short label:** The Bestower of Benefits

**Main description:** Is the One who is almighty to bring inevitably the profit that a creature must receive according to a destiny

**Numbered-file description (unused):** An Naafi', meaning 'The Benefactor' or 'The One Who Causes Benefit', is one of the names attributed to Allah in Islamic tradition. It signifies God's ability to bring goodness, support, and benefit to His creation.

**Displayed summary:** The name An Naafi' highlights the divine attribute of providing benefits and positive outcomes to the believers and all of creation. It emphasizes God's mercy and generosity, reflecting His role as a source of sustenance and care for humanity and the universe.

**Audit:** 30:37 concerns provision rather than a literal name. Prefer 10:107 as a related passage on harm/good, alongside 6:17; explain benefit, paired with Daarr.

- [30:37](https://quran.com/30/37): أَوَلَمۡ يَرَوۡاْ أَنَّ ٱللَّهَ يَبۡسُطُ ٱلرِّزۡقَ لِمَن يَشَآءُ وَيَقۡدِرُۚ إِنَّ فِي ذَٰلِكَ لَأٓيَٰتٖ لِّقَوۡمٖ يُؤۡمِنُونَ

### 93. An Noor — النُّورُ

**Short label:** The Light

**Main description:** The Light: Is the One who guides to their goal those who are prey to temptation and puts them in the right direction according to a destiny; Allah is the One who guides believers to the light of faith, creates light and is not a light that is only a creature

**Numbered-file description (unused):** An Noor, meaning 'The Light', is one of the names of Allah in Islam. It signifies divine illumination, guidance, and the essence of purity and truth. This name emphasizes the importance of spiritual enlightenment and the illumination provided by God in the lives of believers.

**Displayed summary:** The name An Noor originates from the Arabic word 'Noor', which means light. It is associated with the Quranic verse that highlights Allah as the source of light in both the physical and spiritual realms. This name encourages reflection on the guidance and clarity that faith provides to individuals, illuminating their path in life.

**Audit:** 24:35 directly supplies Light. Keep the explanation grounded in the verse and distinguish any interpretive imagery from the text itself.

- [24:35](https://quran.com/24/35): ۞ ٱللَّهُ نُورُ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِۚ مَثَلُ نُورِهِۦ كَمِشۡكَوٰةٖ فِيهَا مِصۡبَاحٌۖ ٱلۡمِصۡبَاحُ فِي زُجَاجَةٍۖ ٱلزُّجَاجَةُ كَأَنَّهَا كَوۡكَبٞ دُرِّيّٞ يُوقَدُ مِن شَجَرَةٖ مُّبَٰرَكَةٖ زَيۡتُونَةٖ لَّا شَرۡقِيَّةٖ وَلَا غَرۡبِيَّةٖ يَكَادُ زَيۡتُهَا يُضِيٓءُ وَلَوۡ لَمۡ تَمۡسَسۡهُ نَارٞۚ نُّورٌ عَلَىٰ نُورٖۚ يَهۡدِي ٱللَّهُ لِنُورِهِۦ مَن يَشَآءُۚ وَيَضۡرِبُ ٱللَّهُ ٱلۡأَمۡثَٰلَ لِلنَّاسِۗ وَٱللَّهُ بِكُلِّ شَيۡءٍ عَلِيمٞ

### 94. Al Haadi — الْهَادِي

**Short label:** The Guider

**Main description:** Is the One who intends for certain creatures the benefit of being well guided and of knowing righteousness

**Numbered-file description (unused):** Al Haadi is one of the names of Allah in Islamic tradition. It translates to 'The Guide' in Arabic.

**Displayed summary:** The name Al Haadi emphasizes Allah's role as the ultimate guide and leader for humanity towards the right path. It signifies His ability to provide direction, wisdom, and understanding, helping individuals make righteous choices in their lives. The concept is rooted in the Quran and reflects the belief that guidance is a divine gift bestowed upon those who seek it sincerely.

**Audit:** 25:31 has haadiyan, a direct nominal form. The Guide is more idiomatic than The Guider.

- [25:31](https://quran.com/25/31): وَكَذَٰلِكَ جَعَلۡنَا لِكُلِّ نَبِيٍّ عَدُوّٗا مِّنَ ٱلۡمُجۡرِمِينَۗ وَكَفَىٰ بِرَبِّكَ هَادِيٗا وَنَصِيرٗا

### 95. Al Badi' — الْبَدِيعُ

**Short label:** The Incomparable Originator

**Main description:** Is the One who brings creatures and this world into existence, not from pre-existing matter or in the image of an earlier model, but on the contrary without precursor, model and unprecedented

**Numbered-file description (unused):** Al Badi' is one of the names of Allah in Islamic tradition, which means 'The Incomparable' or 'The Originator'.

**Displayed summary:** The name Al Badi' signifies Allah's unique ability to create everything from nothing, emphasizing His unmatched creativity and originality in bringing forth the universe. It reflects the concept of being peerless in creation, underscoring that no one can replicate His divine act of creation. This name originates from the Arabic root 'ب-د-ع' (B-D-A), which conveys the sense of initiating and bringing something into existence without any prior example.

**Audit:** 2:117 and 6:101 contain badi' directly. Explain origination without a prior model; overlap with Mubdi is not a data error.

- [2:117](https://quran.com/2/117): بَدِيعُ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِۖ وَإِذَا قَضَىٰٓ أَمۡرٗا فَإِنَّمَا يَقُولُ لَهُۥ كُن فَيَكُونُ
- [6:101](https://quran.com/6/101): بَدِيعُ ٱلسَّمَٰوَٰتِ وَٱلۡأَرۡضِۖ أَنَّىٰ يَكُونُ لَهُۥ وَلَدٞ وَلَمۡ تَكُن لَّهُۥ صَٰحِبَةٞۖ وَخَلَقَ كُلَّ شَيۡءٖۖ وَهُوَ بِكُلِّ شَيۡءٍ عَلِيمٞ

### 96. Al Baaqi — اَلْبَاقِي

**Short label:** The Everlasting One

**Main description:** Is the One whose existence free from annihilation and change is imposed on sound reason which cannot claim the contrary

**Numbered-file description (unused):** Al Baaqi, meaning 'The Ever-Surviving' or 'The Eternal', is one of the names of Allah in Islam. It signifies that Allah is everlasting and exists without an end, in contrast to all created things that are temporary and perishable.

**Displayed summary:** Al-Baaqi signifies that Allah is the sole entity that remains after all of creation ceases to exist. He is the ultimate, unchanging reality, beyond the limitations of time and mortality. While everything else is transient and subject to decay, Allah alone is eternal and enduring, representing the ultimate source of permanence and stability. This name emphasizes His absolute and infinite existence, highlighting that all other things are temporary and dependent on Him.

**Audit:** Remove the unqualified claim that all creation ceases to exist. 55:26–27 contrasts earthly mortality with Allah's remaining; 55:27 uses a verb, not Al Baaqi.

- [55:27](https://quran.com/55/27): وَيَبۡقَىٰ وَجۡهُ رَبِّكَ ذُو ٱلۡجَلَٰلِ وَٱلۡإِكۡرَامِ

### 97. Al Waarith — الْوَارِثُ

**Short label:** The Inheritor

**Main description:** Is the One who exists and does not change as creatures are wiped out

**Numbered-file description (unused):** Al Waarith (الوَارِث) is one of the Names of Allah in Islam, which means 'The Inheritor' or 'The Heir'.

**Displayed summary:** The name Al Waarith signifies that Allah is the ultimate inheritor of all things, as everything in existence belongs to Him, and all creation will return to Him. This concept highlights the transient nature of life and material possessions. The origin of this name is derived from the Quran, specifically in Surah Al-Anbiya (21:55), where Allah is referred to as the Heir of the heavens and the earth. This name instills a sense of reliance on Allah and reflects His eternal nature as the owner and keeper of all that exists.

**Audit:** Correct the summary's false 21:55 attribution. Existing 15:23 has plural al-warithun; 19:40 and 21:89 provide additional relevant passages.

- [15:23](https://quran.com/15/23): وَإِنَّا لَنَحۡنُ نُحۡيِۦ وَنُمِيتُ وَنَحۡنُ ٱلۡوَٰرِثُونَ

### 98. Ar Rasheed — الرَّشِيدُ

**Short label:** The Most Righteous Guide

**Main description:** Is the One who guides creatures to what is in their best interests

**Numbered-file description (unused):** Ar-Rasheed is one of the names of Allah in Islam, which means 'The Guide' or 'The Rightly Guided'. It signifies divine guidance and wisdom from Allah.

**Displayed summary:** The name Ar-Rasheed conveys the idea of Allah being the one who leads and guides believers towards the path of righteousness and truth. The name is derived from the Arabic root 'راشد' (r-sh-d), which pertains to guidance and correctness. It appears in various Islamic teachings, emphasizing the importance of seeking guidance from Allah in all aspects of life.

**Audit:** 2:256 mentions right guidance, not Allah as Ar Rasheed. Correct the summary's claim that راشد is a root: root is ر ش د. Avoid Rightly Guided wording that suggests Allah receives guidance.

- [2:256](https://quran.com/2/256): لَآ إِكۡرَاهَ فِي ٱلدِّينِۖ قَد تَّبَيَّنَ ٱلرُّشۡدُ مِنَ ٱلۡغَيِّۚ فَمَن يَكۡفُرۡ بِٱلطَّٰغُوتِ وَيُؤۡمِنۢ بِٱللَّهِ فَقَدِ ٱسۡتَمۡسَكَ بِٱلۡعُرۡوَةِ ٱلۡوُثۡقَىٰ لَا ٱنفِصَامَ لَهَاۗ وَٱللَّهُ سَمِيعٌ عَلِيمٌ

### 99. As Saboor — الصَّبُورُ

**Short label:** The Patient One

**Main description:** Is the One who delays the punishment of sinners by adjusting the accounts of each until a definite day and gives them their ease and time until a time of which they all do not know the date.

**Numbered-file description (unused):** As Saboor (الْصَّبُور) is one of the names of Allah in Islam, meaning 'The Patient' or 'The Most Patient'. It reflects Allah's attribute of enduring and being patient with His creation, forgiving their shortcomings and sins.

**Displayed summary:** The name As Saboor signifies Allah's enduring patience and ability to withhold anger and punishment, emphasizing His compassion and mercy towards humanity. The concept encourages believers to embody patience in their own lives, understanding the importance of forbearance in facing life's challenges.

**Audit:** Current references command/commend human patience, not a literal divine name. Use 16:61/35:45 as related passages about delaying punishment; Bukhari 6099 supports divine patience, using a comparative form rather than Al Saboor.

- [2:153](https://quran.com/2/153): يَٰٓأَيُّهَا ٱلَّذِينَ ءَامَنُواْ ٱسۡتَعِينُواْ بِٱلصَّبۡرِ وَٱلصَّلَوٰةِۚ إِنَّ ٱللَّهَ مَعَ ٱلصَّٰبِرِينَ
- [3:200](https://quran.com/3/200): يَٰٓأَيُّهَا ٱلَّذِينَ ءَامَنُواْ ٱصۡبِرُواْ وَصَابِرُواْ وَرَابِطُواْ وَٱتَّقُواْ ٱللَّهَ لَعَلَّكُمۡ تُفۡلِحُونَ
- [103:3](https://quran.com/103/3): إِلَّا ٱلَّذِينَ ءَامَنُواْ وَعَمِلُواْ ٱلصَّٰلِحَٰتِ وَتَوَاصَوۡاْ بِٱلۡحَقِّ وَتَوَاصَوۡاْ بِٱلصَّبۡرِ

## Validation result

99 unique sequential IDs; 99 numbered summary files present; 239 listed verse references checked after expanding ranges, all resolve to bundled verses. This verifies reference existence, not semantic suitability. The semantic exceptions are recorded above.

Known linking defect: `82:10-12` becomes `/quran/82/10-12` in both renderers. It needs individual verse links or explicit range support; this audit recommends removing it from Muhsi because its subject is recording angels.

Implementation should update the shared labels/descriptions and the separate summaries together, reconcile or remove unused duplicate descriptions, attach per-reference evidence notes, and invalidate only affected persisted content. Preserve existing fonts, spacing and controls.
