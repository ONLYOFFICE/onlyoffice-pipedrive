import { Tabs, TabList, Tab, TabPanel } from "react-tabs";
import React, { useState } from "react";
import { useSnapshot } from "valtio";
import { Command } from "@pipedrive/app-extensions-sdk";

import { PipedriveSDK } from "@context/PipedriveContext";

import { Creation } from "./Creation";
import { Upload } from "./Upload";

export const CreatePage: React.FC = () => {
  const { sdk } = useSnapshot(PipedriveSDK);
  const [selected, setSelected] = useState(0);
  sdk.execute(Command.RESIZE, {
    height: 500,
    width: 622,
  });

  return (
    <div className="relative w-[622px] h-[500px] flex flex-col overflow-hidden">
      <Tabs
        className="flex justify-center items-start flex-col h-screen"
        onSelect={(index) => setSelected(index)}
      >
        <TabList
          className="flex justify-start items-center min-h-[40px] w-full"
          style={{ backgroundColor: "#F7F7F7" }}
        >
          <Tab
            id="create-file"
            className={`mx-5 outline-none hover:cursor-pointer ${
              selected === 0 ? "text-sky-500" : "text-gray-400"
            }`}
          >
            Create
          </Tab>
          <Tab
            id="upload-file"
            className={`mx-5 outline-none hover:cursor-pointer ${
              selected === 1 ? "text-sky-500" : "text-gray-400"
            }`}
          >
            Upload
          </Tab>
        </TabList>
        <TabPanel
          className={`${
            selected === 0 ? "h-[calc(100%-40px)] w-full" : "h-0 w-0"
          }`}
        >
          <Creation />
        </TabPanel>
        <TabPanel
          className={`${
            selected === 1 ? "h-[calc(100%-40px)] w-full" : "h-0 w-0"
          }`}
        >
          <Upload />
        </TabPanel>
      </Tabs>
    </div>
  );
};