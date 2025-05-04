interface Param {
  name: string;
  type: string;
  description: string;
}

interface ResponseField {
  name: string;
  type: string;
  description: string;
}

interface Endpoint {
  title: string;
  description: string;
  url: string;
  method?: string;
  requestFormat?: string;
  responseFormat: string;
  params?: Param[];
  responseFields?: ResponseField[];
}

const endpoints: Endpoint[] = [
  {
    title: 'Quran',
    description: 'Returns the entire Quran',
    url: '/api/quran',
    responseFormat: 'JSON',
  },
  {
    title: 'List of Quran Chapters',
    description: 'Returns a list of Quran chapters',
    url: '/api/quran/chapters',
    responseFormat: 'JSON',
    responseFields: [
      {
        name: 'name',
        type: 'string',
        description: 'Transliterated name of chapter',
      },
      { name: 'number', type: 'int', description: 'Number of the chapter' },
      {
        name: 'english',
        type: 'string',
        description: 'English name of chapter',
      },
      {
        name: 'verse_count',
        type: 'int',
        description: 'Number of verses in chapter',
      },
    ],
  },
  {
    title: 'Quran by Chapter',
    description: 'Returns a chapter of the quran',
    url: '/api/quran/{chapter}',
    responseFormat: 'JSON',
    responseFields: [
      { name: 'name', type: 'string', description: 'Name of chapter' },
      { name: 'number', type: 'int', description: 'Number of the chapter' },
      { name: 'verses', type: 'array', description: 'Verses in the chapter' },
      { name: 'english', type: 'string', description: 'Name in english' },
    ],
  },
  {
    title: 'Quran by Verse',
    description: 'Returns a verse of the quran',
    url: '/api/quran/{chapter}/{verse}',
    responseFormat: 'JSON',
    responseFields: [
      { name: 'chapter', type: 'int', description: 'Chapter of the verse' },
      { name: 'number', type: 'int', description: 'Number of the verse' },
      { name: 'text', type: 'string', description: 'Text of the verse' },
      {
        name: 'arabic',
        type: 'string',
        description: 'Arabic text of the verse',
      },
    ],
  },
  {
    title: 'Hadith',
    description: 'Returns the entire Hadith',
    url: '/api/hadith',
    responseFormat: 'JSON',
  },
  {
    title: 'Hadith by Book',
    description: 'Returns a book from the hadith',
    url: '/api/hadith/{book}',
    responseFormat: 'JSON',
    responseFields: [
      { name: 'name', type: 'string', description: 'Name of book' },
      { name: 'hadiths', type: 'array', description: 'Hadiths in the book' },
    ],
  },
  {
    title: 'Names',
    description: 'Returns the names of Allah',
    url: '/api/names',
    responseFormat: 'JSON',
  },
  {
    title: 'Search',
    description: 'Get summarised answers via an LLM',
    url: '/api/search',
    method: 'POST',
    requestFormat: 'JSON',
    responseFormat: 'JSON',
    params: [{ name: 'q', type: 'string', description: 'The question to ask' }],
    responseFields: [
      { name: 'q', type: 'string', description: 'The question asked' },
      { name: 'answer', type: 'string', description: 'Answer to the question' },
      {
        name: 'references',
        type: 'array',
        description: 'A list of references used',
      },
    ],
  },
];

export default function Api() {
  return (
    <div className='flex flex-col h-full overflow-y-auto'>
      <div className='max-w-4xl mx-auto px-4 py-4 sm:py-6 lg:py-8 w-full p-0 lg:p-8'>
        <h1 className='text-xl sm:text-2xl md:text-3xl font-bold mb-4 sm:mb-6 lg:mb-8'>Endpoints</h1>
        <p className='mb-4 sm:mb-6 lg:mb-8 text-sm sm:text-base'>A list of API endpoints</p>

        <div className='space-y-6 sm:space-y-8 lg:space-y-12'>
          {endpoints.map((endpoint, index) => (
            <div key={index} className='border-t pt-4 sm:pt-6 lg:pt-8'>
              <h2 className='text-lg sm:text-xl font-semibold mb-2 sm:mb-4'>{endpoint.title}</h2>
              <p className='mb-3 sm:mb-4 text-sm sm:text-base'>{endpoint.description}</p>

              <div className='bg-gray-50 p-2 sm:p-4 font-mono mb-3 sm:mb-4 text-xs sm:text-sm overflow-x-auto'>
                {endpoint.method && (
                  <span className='font-semibold'>{endpoint.method} </span>
                )}
                <span>{endpoint.url}</span>
              </div>

              {endpoint.params && (
                <div className='mb-3 sm:mb-4'>
                  <h3 className='font-semibold mb-1 sm:mb-2 text-sm sm:text-base'>Request Parameters</h3>
                  <p className='mb-2 text-xs sm:text-sm'>Format: {endpoint.requestFormat}</p>
                  <div className='overflow-x-auto'>
                    <table className='w-full border-collapse min-w-full'>
                      <thead>
                        <tr className='border-b'>
                          <th className='text-left py-1 sm:py-2 text-xs sm:text-sm'>Field</th>
                          <th className='text-left py-1 sm:py-2 text-xs sm:text-sm'>Type</th>
                          <th className='text-left py-1 sm:py-2 text-xs sm:text-sm'>Description</th>
                        </tr>
                      </thead>
                      <tbody>
                        {endpoint.params.map((param, idx) => (
                          <tr key={idx} className='border-b'>
                            <td className='py-1 sm:py-2 text-xs sm:text-sm'>{param.name}</td>
                            <td className='py-1 sm:py-2 text-xs sm:text-sm'>{param.type}</td>
                            <td className='py-1 sm:py-2 text-xs sm:text-sm'>{param.description}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              )}

              <div>
                <h3 className='font-semibold mb-1 sm:mb-2 text-sm sm:text-base'>Response</h3>
                <p className='mb-2 text-xs sm:text-sm'>Format: {endpoint.responseFormat}</p>

                {endpoint.responseFields && (
                  <div className='overflow-x-auto'>
                    <table className='w-full border-collapse min-w-full'>
                      <thead>
                        <tr className='border-b'>
                          <th className='text-left py-1 sm:py-2 text-xs sm:text-sm'>Field</th>
                          <th className='text-left py-1 sm:py-2 text-xs sm:text-sm'>Type</th>
                          <th className='text-left py-1 sm:py-2 text-xs sm:text-sm'>Description</th>
                        </tr>
                      </thead>
                      <tbody>
                        {endpoint.responseFields.map((field, idx) => (
                          <tr key={idx} className='border-b'>
                            <td className='py-1 sm:py-2 text-xs sm:text-sm'>{field.name}</td>
                            <td className='py-1 sm:py-2 text-xs sm:text-sm'>{field.type}</td>
                            <td className='py-1 sm:py-2 text-xs sm:text-sm'>{field.description}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
