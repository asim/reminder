import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router';
import { httpGet } from '~/utils/http';

interface LatestResponse {
    name: string;
    hadith: string;
    verse: string;
    links: Record<string, string>;
    updated: string;
    message: string;
}

export function meta() {
    return [
        { title: 'Home - Reminder' },
        {
            property: 'og:title',
            content: 'Home - Reminder',
        },
        {
            name: 'description',
            content: 'The latest reminder from the Quran, Hadith, and Names of Allah',
        },
    ];
}

export default function HomePage() {
    const { data, isLoading, error } = useQuery<LatestResponse>({
        queryKey: ['latest'],
        queryFn: async () => httpGet<LatestResponse>('/api/latest'),
        refetchInterval: 60000, // Refetch every minute
    });

    if (isLoading) return <div className="p-4">Loading...</div>;
    if (error || !data) return <div className="p-4 text-red-500">Failed to load latest reminder.</div>;

    return (
        <div className="h-full w-full overflow-y-auto">
            <div className="max-w-4xl mx-auto w-full p-4 lg:p-8 mb-8 sm:mb-12 space-y-6 sm:space-y-8">
                <div className="text-center">
                    <h1 className="text-2xl sm:text-3xl md:text-4xl font-bold mb-2 sm:mb-3">
                        Latest Reminder
                    </h1>
                    <p className="text-sm sm:text-base text-gray-600 mb-1">
                        Updated hourly with a new verse, hadith, and name
                    </p>
                    <p className="text-xs sm:text-sm text-gray-500">
                        Last updated: {new Date(data.updated).toLocaleString()}
                    </p>
                </div>

                <div className="grid gap-4 sm:gap-6">
                    <section className="bg-blue-50 rounded-lg p-4 sm:p-6 shadow-sm">
                        <div className="flex items-center justify-between mb-3">
                            <h2 className="text-lg sm:text-xl font-semibold text-blue-900">Verse from the Quran</h2>
                            <Link
                                to={data.links['verse']}
                                className="text-xs sm:text-sm text-blue-700 hover:text-blue-900 font-medium"
                            >
                                Read more →
                            </Link>
                        </div>
                        <div className="whitespace-pre-wrap leading-relaxed text-sm sm:text-base text-gray-800">
                            {data.verse}
                        </div>
                    </section>

                    <section className="bg-green-50 rounded-lg p-4 sm:p-6 shadow-sm">
                        <div className="flex items-center justify-between mb-3">
                            <h2 className="text-lg sm:text-xl font-semibold text-green-900">Hadith from Sahih Bukhari</h2>
                            <Link
                                to={data.links['hadith']}
                                className="text-xs sm:text-sm text-green-700 hover:text-green-900 font-medium"
                            >
                                Read more →
                            </Link>
                        </div>
                        <div className="whitespace-pre-wrap leading-relaxed text-sm sm:text-base text-gray-800">
                            {data.hadith}
                        </div>
                    </section>

                    <section className="bg-yellow-50 rounded-lg p-4 sm:p-6 shadow-sm">
                        <div className="flex items-center justify-between mb-3">
                            <h2 className="text-lg sm:text-xl font-semibold text-yellow-900">Name of Allah</h2>
                            <Link
                                to={data.links['name']}
                                className="text-xs sm:text-sm text-yellow-700 hover:text-yellow-900 font-medium"
                            >
                                Read more →
                            </Link>
                        </div>
                        <div className="whitespace-pre-wrap leading-relaxed text-sm sm:text-base text-gray-800">
                            {data.name}
                        </div>
                    </section>
                </div>

                <div className="bg-gray-50 rounded-lg p-4 sm:p-6 border border-gray-200">
                    <h3 className="text-base sm:text-lg font-medium mb-2 text-gray-800">
                        Explore More
                    </h3>
                    <p className="text-sm sm:text-base text-gray-600 mb-4">
                        Browse our full collection or view past daily reminders
                    </p>
                    <div className="flex flex-col sm:flex-row gap-2 sm:gap-3">
                        <Link
                            to="/daily"
                            className="inline-block px-4 py-2 bg-black text-white rounded-md hover:bg-gray-800 transition-colors text-center text-sm sm:text-base"
                        >
                            View Daily Archive
                        </Link>
                        <Link
                            to="/quran"
                            className="inline-block px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 transition-colors text-center text-sm sm:text-base"
                        >
                            Read Quran
                        </Link>
                        <Link
                            to="/hadith"
                            className="inline-block px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 transition-colors text-center text-sm sm:text-base"
                        >
                            Read Hadith
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    );
}
