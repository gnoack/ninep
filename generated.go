package ninep

import (
  "errors"
  "io"
  "log"
)

// size[4] Tauth tag[2] afid[4] uname[s] aname[s]
func writeTauth(w io.Writer, tag uint16, afid uint32, uname string, aname string) error {
  if *debugLog {
    log.Println("->", "Tauth", "tag:", tag, "afid:", afid, "uname:", uname, "aname:", aname)
  }
  size := uint32(4 + 1 + 2 + 4 + (2 + len(uname)) + (2 + len(aname)))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Tauth); err != nil {
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

// size[4] Tauth tag[2] afid[4] uname[s] aname[s]
func readTauth(r io.Reader) (tag uint16, afid uint32, uname string, aname string, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType != Tauth {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &afid); err != nil {
    return
  }
  if err = readString(r, &uname); err != nil {
    return
  }
  if err = readString(r, &aname); err != nil {
    return
  }
  if *debugLog {
    log.Println("->", "Tauth", "tag:", tag, "afid:", afid, "uname:", uname, "aname:", aname)
  }
  return
}

// size[4] Rauth tag[2] aqid[13]
func writeRauth(w io.Writer, tag uint16, aqid Qid) error {
  if *debugLog {
    log.Println("<-", "Rauth", "tag:", tag, "aqid:", aqid)
  }
  size := uint32(4 + 1 + 2 + 13)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Rauth); err != nil {
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

// size[4] Rauth tag[2] aqid[13]
func readRauth(r io.Reader) (aqid Qid, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Rauth {
    err = unexpectedMsgError
    return
  }
  if err = readQid(r, &aqid); err != nil {
    return
  }
  if *debugLog {
    log.Println("<-", "Rauth", "tag:", tag, "aqid:", aqid)
  }
  return
}

// size[4] Tattach tag[2] fid[4] afid[4] uname[s] aname[s]
func writeTattach(w io.Writer, tag uint16, fid uint32, afid uint32, uname string, aname string) error {
  if *debugLog {
    log.Println("->", "Tattach", "tag:", tag, "fid:", fid, "afid:", afid, "uname:", uname, "aname:", aname)
  }
  size := uint32(4 + 1 + 2 + 4 + 4 + (2 + len(uname)) + (2 + len(aname)))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Tattach); err != nil {
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

// size[4] Tattach tag[2] fid[4] afid[4] uname[s] aname[s]
func readTattach(r io.Reader) (tag uint16, fid uint32, afid uint32, uname string, aname string, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType != Tattach {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &fid); err != nil {
    return
  }
  if err = readUint32(r, &afid); err != nil {
    return
  }
  if err = readString(r, &uname); err != nil {
    return
  }
  if err = readString(r, &aname); err != nil {
    return
  }
  if *debugLog {
    log.Println("->", "Tattach", "tag:", tag, "fid:", fid, "afid:", afid, "uname:", uname, "aname:", aname)
  }
  return
}

// size[4] Rattach tag[2] qid[13]
func writeRattach(w io.Writer, tag uint16, qid Qid) error {
  if *debugLog {
    log.Println("<-", "Rattach", "tag:", tag, "qid:", qid)
  }
  size := uint32(4 + 1 + 2 + 13)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Rattach); err != nil {
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

// size[4] Rattach tag[2] qid[13]
func readRattach(r io.Reader) (qid Qid, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Rattach {
    err = unexpectedMsgError
    return
  }
  if err = readQid(r, &qid); err != nil {
    return
  }
  if *debugLog {
    log.Println("<-", "Rattach", "tag:", tag, "qid:", qid)
  }
  return
}

// size[4] Tclunk tag[2] fid[4]
func writeTclunk(w io.Writer, tag uint16, fid uint32) error {
  if *debugLog {
    log.Println("->", "Tclunk", "tag:", tag, "fid:", fid)
  }
  size := uint32(4 + 1 + 2 + 4)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Tclunk); err != nil {
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

// size[4] Tclunk tag[2] fid[4]
func readTclunk(r io.Reader) (tag uint16, fid uint32, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType != Tclunk {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &fid); err != nil {
    return
  }
  if *debugLog {
    log.Println("->", "Tclunk", "tag:", tag, "fid:", fid)
  }
  return
}

// size[4] Rclunk tag[2]
func writeRclunk(w io.Writer, tag uint16) error {
  if *debugLog {
    log.Println("<-", "Rclunk", "tag:", tag)
  }
  size := uint32(4 + 1 + 2)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Rclunk); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  return nil
}

// size[4] Rclunk tag[2]
func readRclunk(r io.Reader) (err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Rclunk {
    err = unexpectedMsgError
    return
  }
  if *debugLog {
    log.Println("<-", "Rclunk", "tag:", tag)
  }
  return
}

// size[4] Rerror tag[2] ename[s]
func writeRerror(w io.Writer, tag uint16, ename string) error {
  if *debugLog {
    log.Println("<-", "Rerror", "tag:", tag, "ename:", ename)
  }
  size := uint32(4 + 1 + 2 + (2 + len(ename)))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Rerror); err != nil {
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

// size[4] Rerror tag[2] ename[s]
func readRerror(r io.Reader) (ename string, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Rerror {
    err = unexpectedMsgError
    return
  }
  if err = readString(r, &ename); err != nil {
    return
  }
  if *debugLog {
    log.Println("<-", "Rerror", "tag:", tag, "ename:", ename)
  }
  return
}

// size[4] Tflush tag[2] oldtag[2]
func writeTflush(w io.Writer, tag uint16, oldtag uint16) error {
  if *debugLog {
    log.Println("->", "Tflush", "tag:", tag, "oldtag:", oldtag)
  }
  size := uint32(4 + 1 + 2 + 2)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Tflush); err != nil {
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

// size[4] Tflush tag[2] oldtag[2]
func readTflush(r io.Reader) (tag uint16, oldtag uint16, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType != Tflush {
    err = unexpectedMsgError
    return
  }
  if err = readUint16(r, &oldtag); err != nil {
    return
  }
  if *debugLog {
    log.Println("->", "Tflush", "tag:", tag, "oldtag:", oldtag)
  }
  return
}

// size[4] Rflush tag[2]
func writeRflush(w io.Writer, tag uint16) error {
  if *debugLog {
    log.Println("<-", "Rflush", "tag:", tag)
  }
  size := uint32(4 + 1 + 2)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Rflush); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  return nil
}

// size[4] Rflush tag[2]
func readRflush(r io.Reader) (err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Rflush {
    err = unexpectedMsgError
    return
  }
  if *debugLog {
    log.Println("<-", "Rflush", "tag:", tag)
  }
  return
}

// size[4] Topen tag[2] fid[4] mode[1]
func writeTopen(w io.Writer, tag uint16, fid uint32, mode uint8) error {
  if *debugLog {
    log.Println("->", "Topen", "tag:", tag, "fid:", fid, "mode:", mode)
  }
  size := uint32(4 + 1 + 2 + 4 + 1)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Topen); err != nil {
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

// size[4] Topen tag[2] fid[4] mode[1]
func readTopen(r io.Reader) (tag uint16, fid uint32, mode uint8, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType != Topen {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &fid); err != nil {
    return
  }
  if err = readUint8(r, &mode); err != nil {
    return
  }
  if *debugLog {
    log.Println("->", "Topen", "tag:", tag, "fid:", fid, "mode:", mode)
  }
  return
}

// size[4] Ropen tag[2] qid[13] iounit[4]
func writeRopen(w io.Writer, tag uint16, qid Qid, iounit uint32) error {
  if *debugLog {
    log.Println("<-", "Ropen", "tag:", tag, "qid:", qid, "iounit:", iounit)
  }
  size := uint32(4 + 1 + 2 + 13 + 4)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Ropen); err != nil {
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

// size[4] Ropen tag[2] qid[13] iounit[4]
func readRopen(r io.Reader) (qid Qid, iounit uint32, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Ropen {
    err = unexpectedMsgError
    return
  }
  if err = readQid(r, &qid); err != nil {
    return
  }
  if err = readUint32(r, &iounit); err != nil {
    return
  }
  if *debugLog {
    log.Println("<-", "Ropen", "tag:", tag, "qid:", qid, "iounit:", iounit)
  }
  return
}

// size[4] Tcreate tag[2] fid[4] name[s] perm[4] mode[1]
func writeTcreate(w io.Writer, tag uint16, fid uint32, name string, perm uint32, mode uint8) error {
  if *debugLog {
    log.Println("->", "Tcreate", "tag:", tag, "fid:", fid, "name:", name, "perm:", perm, "mode:", mode)
  }
  size := uint32(4 + 1 + 2 + 4 + (2 + len(name)) + 4 + 1)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Tcreate); err != nil {
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

// size[4] Tcreate tag[2] fid[4] name[s] perm[4] mode[1]
func readTcreate(r io.Reader) (tag uint16, fid uint32, name string, perm uint32, mode uint8, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType != Tcreate {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &fid); err != nil {
    return
  }
  if err = readString(r, &name); err != nil {
    return
  }
  if err = readUint32(r, &perm); err != nil {
    return
  }
  if err = readUint8(r, &mode); err != nil {
    return
  }
  if *debugLog {
    log.Println("->", "Tcreate", "tag:", tag, "fid:", fid, "name:", name, "perm:", perm, "mode:", mode)
  }
  return
}

// size[4] Rcreate tag[2] qid[13] iounit[4]
func writeRcreate(w io.Writer, tag uint16, qid Qid, iounit uint32) error {
  if *debugLog {
    log.Println("<-", "Rcreate", "tag:", tag, "qid:", qid, "iounit:", iounit)
  }
  size := uint32(4 + 1 + 2 + 13 + 4)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Rcreate); err != nil {
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

// size[4] Rcreate tag[2] qid[13] iounit[4]
func readRcreate(r io.Reader) (qid Qid, iounit uint32, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Rcreate {
    err = unexpectedMsgError
    return
  }
  if err = readQid(r, &qid); err != nil {
    return
  }
  if err = readUint32(r, &iounit); err != nil {
    return
  }
  if *debugLog {
    log.Println("<-", "Rcreate", "tag:", tag, "qid:", qid, "iounit:", iounit)
  }
  return
}

// size[4] Topenfd tag[2] fid[4] mode[1]
func writeTopenfd(w io.Writer, tag uint16, fid uint32, mode uint8) error {
  if *debugLog {
    log.Println("->", "Topenfd", "tag:", tag, "fid:", fid, "mode:", mode)
  }
  size := uint32(4 + 1 + 2 + 4 + 1)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Topenfd); err != nil {
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

// size[4] Topenfd tag[2] fid[4] mode[1]
func readTopenfd(r io.Reader) (tag uint16, fid uint32, mode uint8, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType != Topenfd {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &fid); err != nil {
    return
  }
  if err = readUint8(r, &mode); err != nil {
    return
  }
  if *debugLog {
    log.Println("->", "Topenfd", "tag:", tag, "fid:", fid, "mode:", mode)
  }
  return
}

// size[4] Ropenfd tag[2] qid[13] iounit[4] unixfd[4]
func writeRopenfd(w io.Writer, tag uint16, qid Qid, iounit uint32, unixfd uint32) error {
  if *debugLog {
    log.Println("<-", "Ropenfd", "tag:", tag, "qid:", qid, "iounit:", iounit, "unixfd:", unixfd)
  }
  size := uint32(4 + 1 + 2 + 13 + 4 + 4)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Ropenfd); err != nil {
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

// size[4] Ropenfd tag[2] qid[13] iounit[4] unixfd[4]
func readRopenfd(r io.Reader) (qid Qid, iounit uint32, unixfd uint32, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Ropenfd {
    err = unexpectedMsgError
    return
  }
  if err = readQid(r, &qid); err != nil {
    return
  }
  if err = readUint32(r, &iounit); err != nil {
    return
  }
  if err = readUint32(r, &unixfd); err != nil {
    return
  }
  if *debugLog {
    log.Println("<-", "Ropenfd", "tag:", tag, "qid:", qid, "iounit:", iounit, "unixfd:", unixfd)
  }
  return
}

// size[4] Tread tag[2] fid[4] offset[8] count[4]
func writeTread(w io.Writer, tag uint16, fid uint32, offset uint64, count uint32) error {
  if *debugLog {
    log.Println("->", "Tread", "tag:", tag, "fid:", fid, "offset:", offset, "count:", count)
  }
  size := uint32(4 + 1 + 2 + 4 + 8 + 4)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Tread); err != nil {
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

// size[4] Tread tag[2] fid[4] offset[8] count[4]
func readTread(r io.Reader) (tag uint16, fid uint32, offset uint64, count uint32, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType != Tread {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &fid); err != nil {
    return
  }
  if err = readUint64(r, &offset); err != nil {
    return
  }
  if err = readUint32(r, &count); err != nil {
    return
  }
  if *debugLog {
    log.Println("->", "Tread", "tag:", tag, "fid:", fid, "offset:", offset, "count:", count)
  }
  return
}

// size[4] Rread tag[2] data[count[4]]
func writeRread(w io.Writer, tag uint16, data []byte) error {
  if *debugLog {
    log.Println("<-", "Rread", "tag:", tag, "data:", data)
  }
  size := uint32(4 + 1 + 2 + (4 + len(data)))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Rread); err != nil {
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

// size[4] Rread tag[2] data[count[4]]
func readRread(r io.Reader) (data []byte, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Rread {
    err = unexpectedMsgError
    return
  }
  if err = readByteSlice(r, &data); err != nil {
    return
  }
  if *debugLog {
    log.Println("<-", "Rread", "tag:", tag, "data:", data)
  }
  return
}

// size[4] Twrite tag[2] fid[4] offset[8] data[count[4]]
func writeTwrite(w io.Writer, tag uint16, fid uint32, offset uint64, data []byte) error {
  if *debugLog {
    log.Println("->", "Twrite", "tag:", tag, "fid:", fid, "offset:", offset, "data:", data)
  }
  size := uint32(4 + 1 + 2 + 4 + 8 + (4 + len(data)))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Twrite); err != nil {
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

// size[4] Twrite tag[2] fid[4] offset[8] data[count[4]]
func readTwrite(r io.Reader) (tag uint16, fid uint32, offset uint64, data []byte, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType != Twrite {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &fid); err != nil {
    return
  }
  if err = readUint64(r, &offset); err != nil {
    return
  }
  if err = readByteSlice(r, &data); err != nil {
    return
  }
  if *debugLog {
    log.Println("->", "Twrite", "tag:", tag, "fid:", fid, "offset:", offset, "data:", data)
  }
  return
}

// size[4] Rwrite tag[2] count[4]
func writeRwrite(w io.Writer, tag uint16, count uint32) error {
  if *debugLog {
    log.Println("<-", "Rwrite", "tag:", tag, "count:", count)
  }
  size := uint32(4 + 1 + 2 + 4)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Rwrite); err != nil {
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

// size[4] Rwrite tag[2] count[4]
func readRwrite(r io.Reader) (count uint32, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Rwrite {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &count); err != nil {
    return
  }
  if *debugLog {
    log.Println("<-", "Rwrite", "tag:", tag, "count:", count)
  }
  return
}

// size[4] Tremove tag[2] fid[4]
func writeTremove(w io.Writer, tag uint16, fid uint32) error {
  if *debugLog {
    log.Println("->", "Tremove", "tag:", tag, "fid:", fid)
  }
  size := uint32(4 + 1 + 2 + 4)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Tremove); err != nil {
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

// size[4] Tremove tag[2] fid[4]
func readTremove(r io.Reader) (tag uint16, fid uint32, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType != Tremove {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &fid); err != nil {
    return
  }
  if *debugLog {
    log.Println("->", "Tremove", "tag:", tag, "fid:", fid)
  }
  return
}

// size[4] Rremove tag[2]
func writeRremove(w io.Writer, tag uint16) error {
  if *debugLog {
    log.Println("<-", "Rremove", "tag:", tag)
  }
  size := uint32(4 + 1 + 2)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Rremove); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  return nil
}

// size[4] Rremove tag[2]
func readRremove(r io.Reader) (err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Rremove {
    err = unexpectedMsgError
    return
  }
  if *debugLog {
    log.Println("<-", "Rremove", "tag:", tag)
  }
  return
}

// size[4] Tstat tag[2] fid[4]
func writeTstat(w io.Writer, tag uint16, fid uint32) error {
  if *debugLog {
    log.Println("->", "Tstat", "tag:", tag, "fid:", fid)
  }
  size := uint32(4 + 1 + 2 + 4)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Tstat); err != nil {
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

// size[4] Tstat tag[2] fid[4]
func readTstat(r io.Reader) (tag uint16, fid uint32, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType != Tstat {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &fid); err != nil {
    return
  }
  if *debugLog {
    log.Println("->", "Tstat", "tag:", tag, "fid:", fid)
  }
  return
}

// size[4] Rstat tag[2] stat[n]
func writeRstat(w io.Writer, tag uint16, stat Stat) error {
  if *debugLog {
    log.Println("<-", "Rstat", "tag:", tag, "stat:", stat)
  }
  size := uint32(4 + 1 + 2 + (39 + 8 + len(stat.Name) + len(stat.Uid) + len(stat.Gid) + len(stat.Muid)))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Rstat); err != nil {
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

// size[4] Rstat tag[2] stat[n]
func readRstat(r io.Reader) (stat Stat, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Rstat {
    err = unexpectedMsgError
    return
  }
  // TODO: Why is this doubly size delimited?
  var outerStatSize uint16
  if err = readUint16(r, &outerStatSize); err != nil {
    return
  }
  if err = readStat(r, &stat); err != nil {
    return
  }
  if *debugLog {
    log.Println("<-", "Rstat", "tag:", tag, "stat:", stat)
  }
  return
}

// size[4] Twstat tag[2] fid[4] stat[n]
func writeTwstat(w io.Writer, tag uint16, fid uint32, stat Stat) error {
  if *debugLog {
    log.Println("->", "Twstat", "tag:", tag, "fid:", fid, "stat:", stat)
  }
  size := uint32(4 + 1 + 2 + 4 + (39 + 8 + len(stat.Name) + len(stat.Uid) + len(stat.Gid) + len(stat.Muid)))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Twstat); err != nil {
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

// size[4] Twstat tag[2] fid[4] stat[n]
func readTwstat(r io.Reader) (tag uint16, fid uint32, stat Stat, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType != Twstat {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &fid); err != nil {
    return
  }
  // TODO: Why is this doubly size delimited?
  var outerStatSize uint16
  if err = readUint16(r, &outerStatSize); err != nil {
    return
  }
  if err = readStat(r, &stat); err != nil {
    return
  }
  if *debugLog {
    log.Println("->", "Twstat", "tag:", tag, "fid:", fid, "stat:", stat)
  }
  return
}

// size[4] Rwstat tag[2]
func writeRwstat(w io.Writer, tag uint16) error {
  if *debugLog {
    log.Println("<-", "Rwstat", "tag:", tag)
  }
  size := uint32(4 + 1 + 2)
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Rwstat); err != nil {
    return err
  }
  if err := writeUint16(w, tag); err != nil {
    return err
  }
  return nil
}

// size[4] Rwstat tag[2]
func readRwstat(r io.Reader) (err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Rwstat {
    err = unexpectedMsgError
    return
  }
  if *debugLog {
    log.Println("<-", "Rwstat", "tag:", tag)
  }
  return
}

// size[4] Tversion tag[2] msize[4] version[s]
func writeTversion(w io.Writer, tag uint16, msize uint32, version string) error {
  if *debugLog {
    log.Println("->", "Tversion", "tag:", tag, "msize:", msize, "version:", version)
  }
  size := uint32(4 + 1 + 2 + 4 + (2 + len(version)))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Tversion); err != nil {
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

// size[4] Tversion tag[2] msize[4] version[s]
func readTversion(r io.Reader) (tag uint16, msize uint32, version string, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType != Tversion {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &msize); err != nil {
    return
  }
  if err = readString(r, &version); err != nil {
    return
  }
  if *debugLog {
    log.Println("->", "Tversion", "tag:", tag, "msize:", msize, "version:", version)
  }
  return
}

// size[4] Rversion tag[2] msize[4] version[s]
func writeRversion(w io.Writer, tag uint16, msize uint32, version string) error {
  if *debugLog {
    log.Println("<-", "Rversion", "tag:", tag, "msize:", msize, "version:", version)
  }
  size := uint32(4 + 1 + 2 + 4 + (2 + len(version)))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Rversion); err != nil {
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
func readRversion(r io.Reader) (msize uint32, version string, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Rversion {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &msize); err != nil {
    return
  }
  if err = readString(r, &version); err != nil {
    return
  }
  if *debugLog {
    log.Println("<-", "Rversion", "tag:", tag, "msize:", msize, "version:", version)
  }
  return
}

// size[4] Twalk tag[2] fid[4] newfid[4] nwname*(wname[s])
func writeTwalk(w io.Writer, tag uint16, fid uint32, newfid uint32, nwnames []string) error {
  if *debugLog {
    log.Println("->", "Twalk", "tag:", tag, "fid:", fid, "newfid:", newfid, "nwnames:", nwnames)
  }
  size := uint32(4 + 1 + 2 + 4 + 4 + stringSliceSize(nwnames))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Twalk); err != nil {
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

// size[4] Twalk tag[2] fid[4] newfid[4] nwname*(wname[s])
func readTwalk(r io.Reader) (tag uint16, fid uint32, newfid uint32, nwnames []string, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType != Twalk {
    err = unexpectedMsgError
    return
  }
  if err = readUint32(r, &fid); err != nil {
    return
  }
  if err = readUint32(r, &newfid); err != nil {
    return
  }
  if err = readStringSlice(r, &nwnames); err != nil {
    return
  }
  if *debugLog {
    log.Println("->", "Twalk", "tag:", tag, "fid:", fid, "newfid:", newfid, "nwnames:", nwnames)
  }
  return
}

// size[4] Rwalk tag[2] nwqid*(qid[13])
func writeRwalk(w io.Writer, tag uint16, qids []Qid) error {
  if *debugLog {
    log.Println("<-", "Rwalk", "tag:", tag, "qids:", qids)
  }
  size := uint32(4 + 1 + 2 + (2 + 13*len(qids)))
  if err := writeUint32(w, size); err != nil {
    return err
  }
  if err := writeUint8(w, Rwalk); err != nil {
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

// size[4] Rwalk tag[2] nwqid*(qid[13])
func readRwalk(r io.Reader) (qids []Qid, err error) {
  var size uint32
  if err = readUint32(r, &size); err != nil {
    return
  }
  var msgType uint8
  if err = readUint8(r, &msgType); err != nil {
    return
  }
  var tag uint16
  if err = readUint16(r, &tag); err != nil {
    return
  }
  if msgType == Rerror {
    var errmsg string
    if err = readString(r, &errmsg); err != nil {
      return
    }
    err = errors.New(errmsg)
    return
  }
  if msgType != Rwalk {
    err = unexpectedMsgError
    return
  }
  if err = readQidSlice(r, &qids); err != nil {
    return
  }
  if *debugLog {
    log.Println("<-", "Rwalk", "tag:", tag, "qids:", qids)
  }
  return
}
