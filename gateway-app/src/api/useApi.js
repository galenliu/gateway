import React, {useEffect, useState} from "react";
import axios from 'axios';

function UseApi(url) {

    const [isLoading, setIsLoading] = useState(false)
    const [isError, setIsError] = useState(false)
    const [data, setData] = useState()

    useEffect(() => {
        const fetchData = async () => {
            setIsError(false);
            setIsLoading(true);
            try {
                const result = await axios(url);
                setData(result.data);
            } catch (error) {
                setIsError(true);
            }
            setIsLoading(false);
        };
        fetchData();
    }, [url]);

    return [isLoading, isError, data]
}

function useGetAddonInfo() {

    const [isLoading, setIsLoading] = useState(false)
    const [isError, setIsError] = useState(false)
    const [data, setData] = useState()

    const [isLoading1, isError1, data1] = UseApi("http:127.0.0.1:9090/settings/addons_info")

    useEffect(() => {
        const url = '${data.url}?arch=${ data.arch }&version=${ data.version }&node=${ data.node }&python=${ data.python }'
        const fetchData = async () => {
            setIsError(false);
            setIsLoading(true);
            try {
                const result = await axios(url);
                setData(result.data);
            } catch (error) {
                setIsError(true);
            }
            setIsLoading(false);
        };
        fetchData();
    }, [isLoading1, isError1, data1]);

    return [isLoading, isError, data]
}

export {UseApi, useGetAddonInfo}