select nodules.normalized_nodule_id, series.series_instance_uid, reads.*, measures.*

from
  nodules, series, reads, measures

where
  nodules.series_uid = series.uid
  and nodules.normalized_nodule_id = reads.normalized_nodule_id
  and reads.uid = measures.read_uid
  and measures.nodule_uid = nodules.uid

order by
  series.series_instance_uid, nodules.normalized_nodule_id;
