import React from "react";
import FilterSearch from "../components/common/FilterSearch";

export default function Topics() {
    const [topics, setTopics] = React.useState([]);
    const [searchQuery, setSearchQuery] = React.useState('');

    const [hasNextPage, setHasNextPage] = React.useState(false);
    const [hasPrevPage, setHasPrevPage] = React.useState(false);

    const [endpoint, setEndpoint] = React.useState('/topics?pageSize=10')

    React.useEffect(() => {
        const fetchData = async () => {
            // Fetch topics from the backend
            try {
                const response = await fetch(`http://localhost:3030${endpoint}&search_query=${searchQuery}`);
                const data = await response.json();
                setTopics(data['topics']);
            } catch (error) {
                setTopics([{"error": "Could not fetch topics"}]);
            }
        };
        
        fetchData();
    }, [searchQuery])

    return (
        <div className='topics-page w-full h-auto'>
            <div className='justify-center m-auto flex flex-col gap-3 mt-3'>
                <p className='text-5xl text-center'>Topics</p>
                <FilterSearch setNewQuery={setSearchQuery}/>
                <ul className='list-disc list-inside text-center'>
                    {/* <li><a className='underline' href='#'>IEEE_CC</a></li>
                    <li><a className='underline' href='#'>Google_Toros</a></li>
                    <li><a className='underline' href='#'>Dr_Izaddoost_Club</a></li> */}
                    {
                        topics?.err && <li>{topics.err}</li>
                    }
                    {topics?.err || topics.map((topic, index) => (
                        <li key={index}><a className='underline' href='#'>{topic["name"]}</a></li>
                    ))}
                </ul>
                <div className="flex justify-center mt-4 gap-2">
                    <button
                        className="px-4 py-2 bg-[#860038] hover:bg-[#680018] disabled:bg-gray-500 disabled:hover:cursor-not-allowed text-white rounded-md transition-colors duration-300"
                        onClick={() => setEndpoint(`/topics?pageSize=10&after=${topics[0].id}`)}
                        disabled={!hasPrevPage}
                    >
                        Previous Page
                    </button>
                    <button
                    className="px-4 py-2 bg-[#860038] hover:bg-[#680018] disabled:bg-gray-500 disabled:hover:cursor-not-allowed text-white rounded-md transition-colors duration-300"
                        onClick={() => setEndpoint(`/topics?pageSize=10&before=${topics[topics.length - 1].id}`)}
                        disabled={!hasNextPage} // Disable the button if there is no previous page
                    >
                        Next Page
                    </button>
                </div>
            </div>
        </div>
    );
}
