import React from "react";

import Discussions, { Discussion } from "../components/discussion/Discussions";
import FilterSearch from "../components/common/FilterSearch";

export default function DiscussionsView() {
    return (
        <div>
            <FilterSearch />
            <Discussions>
                <Discussion />
            </Discussions>
        </div>
    );
}