# Quran translations and reading settings

## Editions

- **Mustafa Khattab — The Clear Quran** (default): existing `quran/data/1.json`–`114.json`. Restore the app's substitution of “Most Beneficent” to “Most Compassionate”, including the unnumbered opening and default messages. Published examples: https://quran.com/1/1 and https://quran.com/55/1. This is a targeted reversal of the local substitution, not a claim to have re-audited every word of the existing edition.
- **Muhammad Asad — The Message of the Qur’an**: `quran/data/asad.json`, retrieved 2026-10-18 from https://api.alquran.cloud/v1/quran/en.asad (edition `en.asad`). All 114 chapters and 6,236 verse numbers matched against our existing corpus. The verse strings are preserved verbatim; the upstream API wrapper is replaced with a chapter-to-verse-array map. No generated paraphrases or edits. The unnumbered opening uses this edition's 1:1.

Asad snapshot SHA-256: `6d6c7544e4d83009ed45f5ba2ee3b49881bfc32724516e89199173b065eaf2e3`.

Al Quran Cloud's terms (https://alquran.cloud/terms-and-conditions, retrieved 2026-10-18, page dated 2026-06-14) describe translations as contributed by rights-holders or sourced from public-domain editions and request translator attribution for republication. Asad is attributed in the selector; this does not imply that his work is public domain or covered by the repository's code license. Retain edition attribution with exports.

The API supplies verse text, not Asad's explanatory notes. Settings link to Alim's chapter text and notes: https://www.alim.org/translation/asad/78/. Spot checks of 78:33, 56:22 and 55:56 agree with his interpretation on Alim (minor punctuation differs between digital editions). These are alternative published interpretations; the app does not label one as the sole literal meaning of the Arabic.

Word-level glosses and stored tafsir remain separate datasets. English audio remains Sahih International, read by Ibrahim Walk. Switching text translation does not relabel these as Asad or Khattab. The settings panel states this distinction.

## Behaviour

Reading settings are reached from the cog in the persistent Quran navigation (below the chapter-menu icon on mobile, beside chapter search on desktop) on chapter and verse pages. The modal is a bottom sheet on narrow screens and a compact centred dialog on desktop. It replaces the wrapping checkbox row and includes the display-mode selector. Chapter transcript/share icons sit beside the chapter number under the title; verse actions remain beside their verse. Native dialog behaviour handles focus trapping, Escape and focus restoration.

Text-size controls scale Arabic and English independently. Default sizes, font faces, spacing rules and verse-number styling are preserved. Preferences are remembered locally; explicit URL settings override them for shared links. Old view/word/transliteration preferences migrate. Shared links and embeds include the translation; transcripts use the displayed edition and name it.

Both editions are included in the same chapter API response and thus in the existing background offline download. No third-party request is required to switch editions. The corpus hash changes the offline cache version. If an older saved chapter lacks the requested edition, it displays an explicit update message rather than silently labelling Khattab text as Asad.

Arabic text, word alignment and verse count continue to be checked by the existing integrity test. The additional translation test checks coverage, opening formulas, the source snapshot and removal of the local wording substitution.
