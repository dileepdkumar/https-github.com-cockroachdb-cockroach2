// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

import React from "react";
import moment from "moment";
import { DATE_FORMAT } from "src/util";
import { statisticsClasses } from "../transactionsPageClasses";
import { ISortedTablePagination } from "../../sortedtable";
import { Button } from "src/button";
import { ResultsPerPageLabel } from "src/pagination";

const { statistic, countTitle, lastCleared } = statisticsClasses;

interface TableStatistics {
  pagination: ISortedTablePagination;
  totalCount: number;
  lastReset: Date | string;
  arrayItemName: string;
  activeFilters: number;
  search?: string;
  onClearFilters?: () => void;
}

const renderLastCleared = (lastReset: string | Date) => {
  return `Last cleared ${moment.utc(lastReset).format(DATE_FORMAT)}`;
};

export const TransactionsPageStatistic: React.FC<TableStatistics> = ({
  pagination,
  totalCount,
  lastReset,
  search,
  arrayItemName,
  onClearFilters,
  activeFilters,
}) => {
  const resultsPerPageLabel = (
    <ResultsPerPageLabel
      pagination={{ ...pagination, total: totalCount }}
      pageName={arrayItemName}
      search={search}
    />
  );

  const resultsCountAndClear = (
    <>
      {totalCount} {totalCount === 1 ? "result" : "results"}
      &nbsp;&nbsp;&nbsp;| &nbsp;
      <Button onClick={() => onClearFilters()} type="flat" size="small">
        clear filter
      </Button>
    </>
  );

  return (
    <div className={statistic}>
      <h4 className={countTitle}>
        {activeFilters ? resultsCountAndClear : resultsPerPageLabel}
      </h4>
      <h4 className={lastCleared}>{renderLastCleared(lastReset)}</h4>
    </div>
  );
};
