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
        <div className="h-full w-full overflow-y-auto bg-gradient-to-br from-slate-50 via-white to-slate-50">
            <div className="max-w-5xl mx-auto w-full p-4 lg:p-8 mb-8 sm:mb-12">
                {/* Header */}
                <div className="text-center mb-8 sm:mb-12 pt-4 sm:pt-8">
                    <div className="inline-block mb-4">
                        <div className="w-16 h-1 bg-gradient-to-r from-emerald-500 via-teal-500 to-cyan-500 rounded-full mb-4 mx-auto"></div>
                    </div>
                    <h1 className="text-3xl sm:text-4xl md:text-5xl font-serif font-bold text-gray-900 mb-3 sm:mb-4">
                        Daily Reminder
                    </h1>
                    <p className="text-sm sm:text-base text-gray-600 mb-2 italic">
                        "And remind, for indeed, the reminder benefits the believers" — Quran 51:55
                    </p>
                    <p className="text-xs sm:text-sm text-gray-500 mt-4">
                        Updated hourly • Last refreshed {new Date(data.updated).toLocaleTimeString()}
                    </p>
                </div>

                {/* Content Cards */}
                <div className="space-y-6 sm:space-y-8">
                    {/* Quran Verse */}
                    <article className="group relative bg-white rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-100">
                        <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-blue-500 to-indigo-500"></div>
                        <div className="p-6 sm:p-8">
                            <div className="flex items-start justify-between mb-4">
                                <div>
                                    <div className="text-xs sm:text-sm font-semibold text-blue-600 uppercase tracking-wide mb-1">
                                        Quran
                                    </div>
                                    <h2 className="text-xl sm:text-2xl font-serif font-bold text-gray-900">
                                        Verse from the Holy Book
                                    </h2>
                                </div>
                                <Link
                                    to={data.links['verse']}
                                    className="flex-shrink-0 text-sm text-blue-600 hover:text-blue-700 font-medium opacity-0 group-hover:opacity-100 transition-opacity duration-200 flex items-center gap-1"
                                >
                                    <span>Read context</span>
                                    <span className="text-lg">→</span>
                                </Link>
                            </div>
                            <div className="prose prose-lg max-w-none">
                                <p className="whitespace-pre-wrap leading-relaxed text-gray-700 text-base sm:text-lg font-serif italic">
                                    {data.verse}
                                </p>
                            </div>
                        </div>
                    </article>

                    {/* Hadith */}
                    <article className="group relative bg-white rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-100">
                        <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-emerald-500 to-teal-500"></div>
                        <div className="p-6 sm:p-8">
                            <div className="flex items-start justify-between mb-4">
                                <div>
                                    <div className="text-xs sm:text-sm font-semibold text-emerald-600 uppercase tracking-wide mb-1">
                                        Hadith
                                    </div>
                                    <h2 className="text-xl sm:text-2xl font-serif font-bold text-gray-900">
                                        Prophetic Wisdom
                                    </h2>
                                </div>
                                <Link
                                    to={data.links['hadith']}
                                    className="flex-shrink-0 text-sm text-emerald-600 hover:text-emerald-700 font-medium opacity-0 group-hover:opacity-100 transition-opacity duration-200 flex items-center gap-1"
                                >
                                    <span>Read context</span>
                                    <span className="text-lg">→</span>
                                </Link>
                            </div>
                            <div className="prose prose-lg max-w-none">
                                <p className="whitespace-pre-wrap leading-relaxed text-gray-700 text-base sm:text-lg">
                                    {data.hadith}
                                </p>
                            </div>
                        </div>
                    </article>

                    {/* Name of Allah */}
                    <article className="group relative bg-white rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-100">
                        <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-amber-500 to-orange-500"></div>
                        <div className="p-6 sm:p-8">
                            <div className="flex items-start justify-between mb-4">
                                <div>
                                    <div className="text-xs sm:text-sm font-semibold text-amber-600 uppercase tracking-wide mb-1">
                                        Divine Name
                                    </div>
                                    <h2 className="text-xl sm:text-2xl font-serif font-bold text-gray-900">
                                        Name of Allah
                                    </h2>
                                </div>
                                <Link
                                    to={data.links['name']}
                                    className="flex-shrink-0 text-sm text-amber-600 hover:text-amber-700 font-medium opacity-0 group-hover:opacity-100 transition-opacity duration-200 flex items-center gap-1"
                                >
                                    <span>Learn more</span>
                                    <span className="text-lg">→</span>
                                </Link>
                            </div>
                            <div className="prose prose-lg max-w-none">
                                <p className="whitespace-pre-wrap leading-relaxed text-gray-700 text-base sm:text-lg">
                                    {data.name}
                                </p>
                            </div>
                        </div>
                    </article>
                </div>

                {/* Footer Navigation */}
                <div className="mt-12 sm:mt-16 text-center">
                    <div className="inline-block bg-white rounded-2xl shadow-lg p-6 sm:p-8 border border-gray-100">
                        <h3 className="text-lg sm:text-xl font-serif font-semibold mb-3 text-gray-900">
                            Explore the Collection
                        </h3>
                        <p className="text-sm sm:text-base text-gray-600 mb-6 max-w-md mx-auto">
                            Dive deeper into the Quran, Hadith collections, and the beautiful names of Allah
                        </p>
                        <div className="flex flex-col sm:flex-row gap-3 justify-center">
                            <Link
                                to="/daily"
                                className="inline-flex items-center justify-center px-6 py-3 bg-gradient-to-r from-gray-900 to-gray-800 text-white rounded-lg hover:from-gray-800 hover:to-gray-700 transition-all duration-200 font-medium shadow-md hover:shadow-lg text-sm sm:text-base"
                            >
                                Daily Archive
                            </Link>
                            <Link
                                to="/quran"
                                className="inline-flex items-center justify-center px-6 py-3 bg-white border-2 border-gray-200 rounded-lg hover:border-gray-300 hover:bg-gray-50 transition-all duration-200 font-medium text-sm sm:text-base"
                            >
                                Quran
                            </Link>
                            <Link
                                to="/hadith"
                                className="inline-flex items-center justify-center px-6 py-3 bg-white border-2 border-gray-200 rounded-lg hover:border-gray-300 hover:bg-gray-50 transition-all duration-200 font-medium text-sm sm:text-base"
                            >
                                Hadith
                            </Link>
                            <Link
                                to="/names"
                                className="inline-flex items-center justify-center px-6 py-3 bg-white border-2 border-gray-200 rounded-lg hover:border-gray-300 hover:bg-gray-50 transition-all duration-200 font-medium text-sm sm:text-base"
                            >
                                Names
                            </Link>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
