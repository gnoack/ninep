package ninep

import "io"

// size[4] Tauth tag[2] afid[4] uname[s] aname[s]
func writeTauth(w io.Writer, tag uint16, afid uint32, uname string, aname string) error {
  var size uint32 = 4 + 2 + 2 + 4 + (2 + len(uname)) + (2 + len(aname))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Tauth); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, afid); err != nil {
    return err
  }
  if err := writeString(w, uname); err != nil {
    return err
  }
  if err := writeString(w, aname); err != nil {
    return err
  }
  return nil
}

// size[4] Rauth tag[2] aqid[13]
func writeRauth(w io.Writer, tag uint16, aqid Qid) error {
  var size uint32 = 4 + 2 + 2 + 13
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Rauth); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeQid(w, aqid); err != nil {
    return err
  }
  return nil
}

// size[4] Tattach tag[2] fid[4] afid[4] uname[s] aname[s]
func writeTattach(w io.Writer, tag uint16, fid uint32, afid uint32, uname string, aname string) error {
  var size uint32 = 4 + 2 + 2 + 4 + 4 + (2 + len(uname)) + (2 + len(aname))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Tattach); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, fid); err != nil {
    return err
  }
  if err := writeUint32(w, afid); err != nil {
    return err
  }
  if err := writeString(w, uname); err != nil {
    return err
  }
  if err := writeString(w, aname); err != nil {
    return err
  }
  return nil
}

// size[4] Rattach tag[2] qid[13]
func writeRattach(w io.Writer, tag uint16, qid Qid) error {
  var size uint32 = 4 + 2 + 2 + 13
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Rattach); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeQid(w, qid); err != nil {
    return err
  }
  return nil
}

// size[4] Tclunk tag[2] fid[4]
func writeTclunk(w io.Writer, tag uint16, fid uint32) error {
  var size uint32 = 4 + 2 + 2 + 4
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Tclunk); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, fid); err != nil {
    return err
  }
  return nil
}

// size[4] Rclunk tag[2]
func writeRclunk(w io.Writer, tag uint16) error {
  var size uint32 = 4 + 2 + 2
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Rclunk); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  return nil
}

// size[4] Rerror tag[2] ename[s]
func writeRerror(w io.Writer, tag uint16, ename string) error {
  var size uint32 = 4 + 2 + 2 + (2 + len(ename))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Rerror); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeString(w, ename); err != nil {
    return err
  }
  return nil
}

// size[4] Tflush tag[2] oldtag[2]
func writeTflush(w io.Writer, tag uint16, oldtag uint16) error {
  var size uint32 = 4 + 2 + 2 + 2
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Tflush); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint16(w, oldtag); err != nil {
    return err
  }
  return nil
}

// size[4] Rflush tag[2]
func writeRflush(w io.Writer, tag uint16) error {
  var size uint32 = 4 + 2 + 2
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Rflush); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  return nil
}

// size[4] Topen tag[2] fid[4] mode[1]
func writeTopen(w io.Writer, tag uint16, fid uint32, mode uint8) error {
  var size uint32 = 4 + 2 + 2 + 4 + 1
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Topen); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, fid); err != nil {
    return err
  }
  if err := writeUint8(w, mode); err != nil {
    return err
  }
  return nil
}

// size[4] Ropen tag[2] qid[13] iounit[4]
func writeRopen(w io.Writer, tag uint16, qid Qid, iounit uint32) error {
  var size uint32 = 4 + 2 + 2 + 13 + 4
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Ropen); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeQid(w, qid); err != nil {
    return err
  }
  if err := writeUint32(w, iounit); err != nil {
    return err
  }
  return nil
}

// size[4] Tcreate tag[2] fid[4] name[s] perm[4] mode[1]
func writeTcreate(w io.Writer, tag uint16, fid uint32, name string, perm uint32, mode uint8) error {
  var size uint32 = 4 + 2 + 2 + 4 + (2 + len(name)) + 4 + 1
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Tcreate); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, fid); err != nil {
    return err
  }
  if err := writeString(w, name); err != nil {
    return err
  }
  if err := writeUint32(w, perm); err != nil {
    return err
  }
  if err := writeUint8(w, mode); err != nil {
    return err
  }
  return nil
}

// size[4] Rcreate tag[2] qid[13] iounit[4]
func writeRcreate(w io.Writer, tag uint16, qid Qid, iounit uint32) error {
  var size uint32 = 4 + 2 + 2 + 13 + 4
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Rcreate); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeQid(w, qid); err != nil {
    return err
  }
  if err := writeUint32(w, iounit); err != nil {
    return err
  }
  return nil
}

// size[4] Topenfd tag[2] fid[4] mode[1]
func writeTopenfd(w io.Writer, tag uint16, fid uint32, mode uint8) error {
  var size uint32 = 4 + 2 + 2 + 4 + 1
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Topenfd); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, fid); err != nil {
    return err
  }
  if err := writeUint8(w, mode); err != nil {
    return err
  }
  return nil
}

// size[4] Ropenfd tag[2] qid[13] iounit[4] unixfd[4]
func writeRopenfd(w io.Writer, tag uint16, qid Qid, iounit uint32, unixfd uint32) error {
  var size uint32 = 4 + 2 + 2 + 13 + 4 + 4
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Ropenfd); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeQid(w, qid); err != nil {
    return err
  }
  if err := writeUint32(w, iounit); err != nil {
    return err
  }
  if err := writeUint32(w, unixfd); err != nil {
    return err
  }
  return nil
}

// size[4] Tread tag[2] fid[4] offset[8] count[4]
func writeTread(w io.Writer, tag uint16, fid uint32, offset uint64, count uint32) error {
  var size uint32 = 4 + 2 + 2 + 4 + 8 + 4
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Tread); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, fid); err != nil {
    return err
  }
  if err := writeUint64(w, offset); err != nil {
    return err
  }
  if err := writeUint32(w, count); err != nil {
    return err
  }
  return nil
}

// size[4] Rread tag[2] data[count[4]]
func writeRread(w io.Writer, tag uint16, data []byte) error {
  var size uint32 = 4 + 2 + 2 + (4 + len(data))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Rread); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeByteSlice(w, data); err != nil {
    return err
  }
  return nil
}

// size[4] Twrite tag[2] fid[4] offset[8] data[count[4]]
func writeTwrite(w io.Writer, tag uint16, fid uint32, offset uint64, data []byte) error {
  var size uint32 = 4 + 2 + 2 + 4 + 8 + (4 + len(data))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Twrite); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, fid); err != nil {
    return err
  }
  if err := writeUint64(w, offset); err != nil {
    return err
  }
  if err := writeByteSlice(w, data); err != nil {
    return err
  }
  return nil
}

// size[4] Rwrite tag[2] count[4]
func writeRwrite(w io.Writer, tag uint16, count uint32) error {
  var size uint32 = 4 + 2 + 2 + 4
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Rwrite); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, count); err != nil {
    return err
  }
  return nil
}

// size[4] Tremove tag[2] fid[4]
func writeTremove(w io.Writer, tag uint16, fid uint32) error {
  var size uint32 = 4 + 2 + 2 + 4
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Tremove); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, fid); err != nil {
    return err
  }
  return nil
}

// size[4] Rremove tag[2]
func writeRremove(w io.Writer, tag uint16) error {
  var size uint32 = 4 + 2 + 2
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Rremove); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  return nil
}

// size[4] Tstat tag[2] fid[4]
func writeTstat(w io.Writer, tag uint16, fid uint32) error {
  var size uint32 = 4 + 2 + 2 + 4
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Tstat); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, fid); err != nil {
    return err
  }
  return nil
}

// size[4] Rstat tag[2] stat[n]
func writeRstat(w io.Writer, tag uint16, stat Stat) error {
  var size uint32 = 4 + 2 + 2 + (39 + 8 + len(stat.Name) + len(stat.Uid) + len(stat.Gid) + len(stat.Muid))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Rstat); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeStat(w, stat); err != nil {
    return err
  }
  return nil
}

// size[4] Twstat tag[2] fid[4] stat[n]
func writeTwstat(w io.Writer, tag uint16, fid uint32, stat Stat) error {
  var size uint32 = 4 + 2 + 2 + 4 + (39 + 8 + len(stat.Name) + len(stat.Uid) + len(stat.Gid) + len(stat.Muid))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Twstat); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, fid); err != nil {
    return err
  }
  if err := writeStat(w, stat); err != nil {
    return err
  }
  return nil
}

// size[4] Rwstat tag[2]
func writeRwstat(w io.Writer, tag uint16) error {
  var size uint32 = 4 + 2 + 2
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Rwstat); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  return nil
}

// size[4] Tversion tag[2] msize[4] version[s]
func writeTversion(w io.Writer, tag uint16, msize uint32, version string) error {
  var size uint32 = 4 + 2 + 2 + 4 + (2 + len(version))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Tversion); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, msize); err != nil {
    return err
  }
  if err := writeString(w, version); err != nil {
    return err
  }
  return nil
}

// size[4] Rversion tag[2] msize[4] version[s]
func writeRversion(w io.Writer, tag uint16, msize uint32, version string) error {
  var size uint32 = 4 + 2 + 2 + 4 + (2 + len(version))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Rversion); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, msize); err != nil {
    return err
  }
  if err := writeString(w, version); err != nil {
    return err
  }
  return nil
}

// size[4] Twalk tag[2] fid[4] newfid[4] nwname*(wname[s])
func writeTwalk(w io.Writer, tag uint16, fid uint32, newfid uint32, nwnames []string) error {
  var size uint32 = 4 + 2 + 2 + 4 + 4 + stringSliceSize(nwnames)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Twalk); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeUint32(w, fid); err != nil {
    return err
  }
  if err := writeUint32(w, newfid); err != nil {
    return err
  }
  if err := writeStringSlice(w, nwnames); err != nil {
    return err
  }
  return nil
}

// size[4] Rwalk tag[2] nwqid*(qid[13])
func writeRwalk(w io.Writer, tag uint16, qids []Qid) error {
  var size uint32 = 4 + 2 + 2 + (2 + 13*len(qids))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint16(w, Rwalk); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  if err := writeQidSlice(w, qids); err != nil {
    return err
  }
  return nil
}
